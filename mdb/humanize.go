package mdb

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
)

// IDToHuman returns a human readable ID version of a ObjectID
// id	: the ObjectID to convert
func IDToHuman(id bson.ObjectId) (text string) {
	return fmt.Sprintf(`%x`, string(id))
}

// HumanToID returns an ObjectID from a human readable ID
// text	: the human readable ID
func HumanToID(text string) (id bson.ObjectId) {
	if bson.IsObjectIdHex(text) {
		return bson.ObjectIdHex(text)
	}
	return id
}
