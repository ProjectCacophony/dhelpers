package mdb

import (
	"github.com/globalsign/mgo"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// Collection is the Type used for MongoDB Collections
type Collection string

// DB returns the MongoDB Session Collection of a MongoDB Collection
func (c Collection) DB() *mgo.Collection {
	return cache.GetMgo().C(string(c))
}
