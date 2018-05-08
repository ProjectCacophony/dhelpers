package dhelpers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo/bson"
	"github.com/kennygrant/sanitize"
	"github.com/minio/minio-go"
	"github.com/satori/go.uuid"
	"gitlab.com/Cacophony/dhelpers/cache"
	"gitlab.com/Cacophony/dhelpers/mdb"
	"gitlab.com/Cacophony/dhelpers/models"
	"gitlab.com/Cacophony/dhelpers/state"
)

// TODO: watch cache folder size

// AddFileMetadata defines possible metadta for new objects
type AddFileMetadata struct {
	Filename           string            // the actual file name, can be empty
	ChannelID          string            // the source channel ID, can be empty, but should be set if possible
	UserID             string            // the source user ID, can be empty, but should be set if possible
	GuildID            string            // the source guild ID, can be empty but should be set if possible, will be set automatically if ChannelID has been set
	AdditionalMetadata map[string]string // additional metadata attached to the object
}

// AddFile stores a file
// name		: the name of the new object, can be empty to generate an unique name
// data		: the file data
// metadata	: metadata attached to the object
// source	: the source name for the file, for example the module name, can not be empty
// public	: if true file will be available via the website proxy
// TODO: prevent duplicates
func AddFile(name string, data []byte, metadata AddFileMetadata, source string, public bool) (objectName string, err error) {
	// check if source is set
	if source == "" {
		return "", errors.New("source can not be empty")
	}
	// TODO: check if user is allowed to upload files
	// set new object name
	objectName = name
	if objectName == "" {
		// generate unique filename
		objectName = uuid.NewV4().String()
	}
	// retrieve guildID if channelID is set
	guildID := metadata.GuildID
	if metadata.ChannelID != "" {
		var channel *discordgo.Channel
		channel, err = state.Channel(metadata.ChannelID)
		LogError(err)
		if err == nil {
			guildID = channel.GuildID
		}
	}
	// get filetype
	filetype := http.DetectContentType(data)
	// get filesize
	filesize := binary.Size(data)
	// update metadata
	if metadata.AdditionalMetadata == nil {
		metadata.AdditionalMetadata = make(map[string]string)
	}
	metadata.AdditionalMetadata["filename"] = metadata.Filename
	metadata.AdditionalMetadata["userid"] = metadata.UserID
	metadata.AdditionalMetadata["guildid"] = guildID
	metadata.AdditionalMetadata["channelid"] = metadata.ChannelID
	metadata.AdditionalMetadata["source"] = source
	metadata.AdditionalMetadata["mimetype"] = filetype
	metadata.AdditionalMetadata["filesize"] = strconv.Itoa(filesize)
	metadata.AdditionalMetadata["public"] = "no"
	if public {
		metadata.AdditionalMetadata["public"] = "yes"
	}
	// upload file
	err = uploadFile(objectName, data, metadata.AdditionalMetadata)
	if err != nil {
		return "", err
	}
	// store in database
	err = mdb.UpsertQuery(
		models.StorageTable,
		bson.M{"objectname": objectName},
		models.StorageEntry{
			ObjectName:     objectName,
			ObjectNameHash: GetMD5Hash(objectName),
			UploadDate:     time.Now(),
			Filename:       metadata.Filename,
			UserID:         metadata.UserID,
			GuildID:        guildID,
			ChannelID:      metadata.ChannelID,
			Source:         source,
			MimeType:       filetype,
			Filesize:       filesize,
			Public:         public,
			Metadata:       metadata.AdditionalMetadata,
		},
	)
	if err != nil {
		return "", err
	}
	// TODO: warm up cache for public files
	cache.GetLogger().WithField("module", "storage").Infof(
		"stored #%s for %s (%+v)",
		objectName, source, metadata,
	)
	// return new objectName
	return objectName, nil
}

// RetrieveFileInformation retrieves information about a file
// objectName	: the name of the file to retrieve
func RetrieveFileInformation(objectName string) (info models.StorageEntry, err error) {
	err = mdb.One(
		models.StorageTable.DB().Find(bson.M{"objectname": objectName}),
		&info,
	)
	return info, err
}

// RetrieveFile retrieves a file
// objectName	: the name of the file to retrieve
func RetrieveFile(objectName string) (data []byte, err error) {
	// Increase MongoDB RetrievedCount
	go func() {
		defer RecoverLog()
		goErr := mdb.UpdateQuery(models.StorageTable, bson.M{"objectname": objectName}, bson.M{"$inc": bson.M{"retrievedcount": 1}})
		if goErr != nil && !mdb.ErrNotFound(goErr) {
			CheckErr(goErr)
		}
	}()

	data = getBucketCache(objectName)
	if data != nil {
		cache.GetLogger().WithField("module", "storage").Infof("retrieving " + objectName + " from minio cache")
		return data, nil
	}

	cache.GetLogger().WithField("module", "storage").Infof("retrieving " + objectName + " from minio storage")

	// retrieve the object
	minioObject, err := cache.GetMinio().GetObject(getBucket(), sanitize.BaseName(objectName), minio.GetObjectOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "Please reduce your request rate.") {
			cache.GetLogger().WithField("module", "storage").Infof("object storage ratelimited, waiting for one second, then retrying")
			time.Sleep(1 * time.Second)
			return RetrieveFile(objectName)
		}
		if strings.Contains(err.Error(), "net/http") || strings.Contains(err.Error(), "timeout") {
			cache.GetLogger().WithField("module", "storage").Infof("network error retrieving, waiting for one second, then retrying")
			time.Sleep(1 * time.Second)
			return RetrieveFile(objectName)
		}
		return data, err
	}

	// read the object into a byte slice
	data, err = ioutil.ReadAll(minioObject)
	if err != nil {
		return data, err
	}

	go func() {
		defer RecoverLog()
		cache.GetLogger().WithField("module", "storage").Infof("caching " + objectName + " into minio cache")
		err := setBucketCache(objectName, data)
		CheckErr(err)
	}()

	return data, nil
}

// RetrieveFileByHash retrieves a file by the object name md5 hash
// hash	: the md5 hash
func RetrieveFileByHash(hash string) (filename, filetype string, data []byte, err error) {
	var entryBucket models.StorageEntry
	err = mdb.One(
		models.StorageTable.DB().Find(bson.M{"objectnamehash": hash}),
		&entryBucket,
	)
	if err != nil && !mdb.ErrNotFound(err) {
		return "", "", nil, err
	}

	data, err = RetrieveFile(entryBucket.ObjectName)
	if err != nil {
		return "", "", nil, err
	}
	return entryBucket.Filename, entryBucket.MimeType, data, nil
}

// RetrieveFilesByAdditionalObjectMetadata retrieves files by additional object metadta
// currently supported file sources: custom commands
// hash	: the md5 hash
func RetrieveFilesByAdditionalObjectMetadata(key, value string) (objectNames []string, err error) {
	var entryBucket []models.StorageEntry
	err = mdb.Iter(models.StorageTable.DB().Find(
		bson.M{"metadata." + strings.ToLower(key): value},
	)).All(&entryBucket)
	if err != nil {
		return nil, err
	}

	objectNames = make([]string, 0)
	if len(entryBucket) > 0 {
		for _, entry := range entryBucket {
			objectNames = append(objectNames, entry.ObjectName)
		}
	}

	if len(objectNames) < 1 {
		return nil, errors.New("none matching files found")
	}

	return objectNames, nil
}

// DeleteFile deletes a file
// objectName	: the name of the object
func DeleteFile(objectName string) (err error) {
	cache.GetLogger().WithField("module", "storage").Infof("deleting " + objectName + " from minio storage")

	go func() {
		defer RecoverLog()
		cache.GetLogger().WithField("module", "storage").Infof("deleting " + objectName + " from minio cache")
		goErr := deleteBucketCache(objectName)
		CheckErr(goErr)
	}()

	// delete the object
	err = cache.GetMinio().RemoveObject(getBucket(), sanitize.BaseName(objectName))

	// delete mongo db entry
	go func() {
		defer RecoverLog()
		goErr := mdb.DeleteQuery(models.StorageTable, bson.M{"objectname": objectName})
		if goErr != nil && !mdb.ErrNotFound(goErr) {
			CheckErr(err)
		}
	}()

	return err
}

// uploads a file to the minio object storage
// objectName	: the name of the file to upload
// data			: the data for the new object
// metadata		: additional metadata attached to the object
// TODO: prevent overwrites
func uploadFile(objectName string, data []byte, metadata map[string]string) (err error) {
	options := minio.PutObjectOptions{}

	// add content type
	options.ContentType = http.DetectContentType(data)

	// add metadata
	if len(metadata) > 0 {
		options.UserMetadata = metadata
	}

	// upload the data
	_, err = cache.GetMinio().PutObject(getBucket(), sanitize.BaseName(objectName), bytes.NewReader(data), -1, options)
	return err
}

func getBucketCache(objectName string) (data []byte) {
	var err error

	if _, err = os.Stat(getObjectPath(objectName)); os.IsNotExist(err) {
		return nil
	}

	data, err = ioutil.ReadFile(getObjectPath(objectName))
	if err != nil {
		return nil
	}

	return data
}

func setBucketCache(objectName string, data []byte) (err error) {
	if _, err = os.Stat(filepath.Dir(getObjectPath(objectName))); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(getObjectPath(objectName)), os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(getObjectPath(objectName), data, 0644)
	return err
}

func deleteBucketCache(objectName string) (err error) {
	if _, err = os.Stat(getObjectPath(objectName)); os.IsNotExist(err) {
		return nil
	}

	err = os.Remove(getObjectPath(objectName))
	return err
}

func getObjectPath(objectName string) (path string) {
	return os.Getenv("S3_CACHE_FOLDER") + "/" + sanitize.BaseName(objectName)
}

func getBucket() (bucket string) {
	return os.Getenv("S3_BUCKET")
}
