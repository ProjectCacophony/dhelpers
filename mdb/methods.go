package mdb

import (
	"errors"
	"reflect"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Insert inserts data into the database, the struct has to have a field called ID
func Insert(collection Collection, data interface{}) (rid bson.ObjectId, err error) {
	var recordData reflect.Value
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		// handle non pointers
		recordData = reflect.New(reflect.TypeOf(data)).Elem()
		recordData.Set(reflect.ValueOf(data))
	} else {
		// handle pointers
		// convert the raw interface data to its actual type
		recordData = reflect.ValueOf(data).Elem()
	}

	// confirm data has an ID field
	idField := recordData.FieldByName("ID")
	if !idField.IsValid() {
		return bson.ObjectId(""), errors.New("invalid data")
	}

	// if the records id field isn't empty, give it an id
	newID := idField.String()
	if newID == "" {
		newID = string(bson.NewObjectId())
		idField.SetString(newID)
	}

	err = collection.DB().Insert(recordData.Interface())
	if err != nil {
		return bson.ObjectId(""), err
	}

	return bson.ObjectId(newID), nil
}

// UpdateID updates an entry by ID
func UpdateID(collection Collection, id bson.ObjectId, data interface{}) (err error) {
	if !id.Valid() {
		return errors.New("invalid id")
	}

	err = collection.DB().UpdateId(id, data)
	return err
}

// UpdateQuery updates an entry by query
func UpdateQuery(collection Collection, selector interface{}, data interface{}) (err error) {
	err = collection.DB().Update(selector, data)
	return err
}

// UpsertID upserts an entry by ID
func UpsertID(collection Collection, id bson.ObjectId, data interface{}) (err error) {
	if !id.Valid() {
		return errors.New("invalid id")
	}

	_, err = collection.DB().UpsertId(id, data)
	return err
}

// UpsertQuery upserts an entry by query
func UpsertQuery(collection Collection, selector interface{}, data interface{}) (err error) {
	_, err = collection.DB().Upsert(selector, data)
	return err
}

// DeleteID deletes an entry by ID
func DeleteID(collection Collection, id bson.ObjectId) (err error) {
	if !id.Valid() {
		return errors.New("invalid id")
	}

	err = collection.DB().RemoveId(id)
	return err
}

// DeleteQuery deletes an entry by query
func DeleteQuery(collection Collection, selector interface{}) (err error) {
	err = collection.DB().Remove(selector)
	return err
}

// Iter starts an iteration
func Iter(query *mgo.Query) (iter *mgo.Iter) {
	iter = query.Iter()
	return iter
}

// One unmarshalls one query result into object
func One(query *mgo.Query, object interface{}) (err error) {
	err = query.One(object)
	return err
}

// PipeOne starts a pipeline
func PipeOne(collection Collection, pipeline interface{}, object interface{}) (err error) {
	err = collection.DB().Pipe(pipeline).One(object)
	return err
}

// Count counts all entries matching a query
func Count(collection Collection, query interface{}) (count int, err error) {
	count, err = collection.DB().Find(query).Count()
	return count, err
}
