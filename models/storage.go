package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.com/project-d-collab/dhelpers/mdb"
)

const (
	// StorageTable is the table containing all StorageEntry entries
	StorageTable mdb.Collection = "storage"
)

// StorageEntry contains information about an object stored in object storage
type StorageEntry struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
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
