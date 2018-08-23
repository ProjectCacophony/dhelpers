package mongo

import (
	"context"

	"reflect"

	"errors"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// Collection is the type for Database Collections
type Collection string

// BasicRepository is a Repository with common database logic
type BasicRepository interface {
	GetByID(ctx context.Context, id objectid.ObjectID, result interface{}) error
	Find(ctx context.Context, filter interface{}, result interface{}) error
	FindOne(ctx context.Context, filter interface{}, result interface{}) error
	UpdateByID(ctx context.Context, id objectid.ObjectID, document interface{}) error
	Update(ctx context.Context, filter interface{}, document interface{}) error
	UpsertByID(ctx context.Context, id objectid.ObjectID, document interface{}) error
	Upsert(ctx context.Context, filter interface{}, document interface{}) error
	Store(ctx context.Context, document interface{}) (*objectid.ObjectID, error)
	DeleteByID(ctx context.Context, id objectid.ObjectID) error
	Delete(ctx context.Context, filter interface{}) error
	Count(ctx context.Context, filter interface{}) (int64, error)
}

// NewRepository creates a new MongoDB Repository from a MongoDB collection with the BasicRepository type
func NewRepository(collection Collection) BasicRepository {
	return &basicRepositoryUsecase{
		collectionName: collection,
	}
}

type basicRepositoryUsecase struct {
	collectionName Collection
	collection     *mongo.Collection
}

// TODO: tracing

// "lazy loading" collection because mongo DB might not have been initialised yet
func (r *basicRepositoryUsecase) initCollection() error {
	if r.collection == nil {
		if cache.GetMongo() == nil {
			return ErrUnavailable
		}

		r.collection = cache.GetMongo().Collection(string(r.collectionName))
	}
	return nil
}

func (r *basicRepositoryUsecase) GetByID(ctx context.Context, id objectid.ObjectID, result interface{}) error {
	return r.FindOne(ctx, map[string]objectid.ObjectID{"_id": id}, &result)
}

// based on https://github.com/globalsign/mgo/blob/master/session.go#L4428
func (r *basicRepositoryUsecase) Find(ctx context.Context, filter interface{}, result interface{}) error {
	err := r.initCollection()
	if err != nil {
		return err
	}

	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr {
		return errors.New("result argument must be a slice address")
	}

	slicev := resultv.Elem()

	if slicev.Kind() == reflect.Interface {
		slicev = slicev.Elem()
	}
	if slicev.Kind() != reflect.Slice {
		return errors.New("result argument must be a slice address")
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			cache.GetLogger().WithError(err).Errorln("error closing cursor")
		}
	}()

	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	var i int
	for cursor.Next(ctx) {
		if slicev.Len() == i {
			elemp := reflect.New(elemt)

			err := cursor.Decode(elemp.Interface())
			if err != nil {
				return err
			}

			slicev = reflect.Append(slicev, elemp.Elem())
			slicev = slicev.Slice(0, slicev.Cap())
		} else {
			err := cursor.Decode(slicev.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}

		i++
	}
	resultv.Elem().Set(slicev.Slice(0, i))

	return cursor.Err()
}

func (r *basicRepositoryUsecase) FindOne(ctx context.Context, filter interface{}, document interface{}) error {
	err := r.initCollection()
	if err != nil {
		return err
	}

	docResult := r.collection.FindOne(ctx, filter)
	if docResult == nil {
		return ErrNotFound
	}

	err = docResult.Decode(document)
	if err == mongo.ErrNoDocuments {
		return ErrNotFound
	}
	return err
}

func (r *basicRepositoryUsecase) UpdateByID(ctx context.Context, id objectid.ObjectID, document interface{}) error {
	return r.Update(ctx, map[string]objectid.ObjectID{"_id": id}, document)
}

func (r *basicRepositoryUsecase) Update(ctx context.Context, filter interface{}, document interface{}) error {
	err := r.initCollection()
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(ctx, filter, document)
	if err != nil {
		return err
	}

	if result.MatchedCount <= 0 || result.ModifiedCount <= 0 {
		return ErrNotFound
	}
	return nil
}

func (r *basicRepositoryUsecase) UpsertByID(ctx context.Context, id objectid.ObjectID, document interface{}) error {
	return r.Upsert(ctx, map[string]objectid.ObjectID{"_id": id}, document)
}

func (r *basicRepositoryUsecase) Upsert(ctx context.Context, filter interface{}, document interface{}) error {
	err := r.initCollection()
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(ctx, filter, document, updateopt.Upsert(true))
	return err
}

func (r *basicRepositoryUsecase) Store(ctx context.Context, document interface{}) (*objectid.ObjectID, error) {
	err := r.initCollection()
	if err != nil {
		return nil, err
	}

	result, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	id, ok := result.InsertedID.(objectid.ObjectID)
	if !ok {
		return nil, errors.New("error gathering object ID")
	}
	return &id, nil
}

func (r *basicRepositoryUsecase) DeleteByID(ctx context.Context, id objectid.ObjectID) error {
	return r.Delete(ctx, map[string]objectid.ObjectID{"_id": id})
}

func (r *basicRepositoryUsecase) Delete(ctx context.Context, filter interface{}) error {
	err := r.initCollection()
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount <= 0 {
		return ErrNotFound
	}
	return nil
}

func (r *basicRepositoryUsecase) Count(ctx context.Context, filter interface{}) (int64, error) {
	err := r.initCollection()
	if err != nil {
		return 0, err
	}

	return r.collection.Count(ctx, filter)
}
