package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"gitlab.com/Cacophony/dhelpers/mongo"
)

const (
	// StorageTable is the table containing all StorageEntry entries
	StorageTable mongo.Collection = "storage"
)

var (
	// StorageRepository contains the database logic for the table
	StorageRepository = mongo.NewRepository(StorageTable)
)

// StorageEntry contains information about an object stored in object storage
type StorageEntry struct {
	ID             *objectid.ObjectID `bson:"_id,omitempty"`
	ObjectName     string
	ObjectNameHash string
	UploadDate     time.Time
	Filename       string
	UserID         string
	GuildID        string
	ChannelID      string
	Source         string
	MimeType       string
	Filesize       int // in bytes
	Public         bool
	Metadata       map[string]string
	RetrievedCount int
}
