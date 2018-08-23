package components

import (
	"os"

	"time"

	"context"

	"github.com/mongodb/mongo-go-driver/core/readconcern"
	"github.com/mongodb/mongo-go-driver/core/writeconcern"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// InitMongoDB initialises the MongoDB session
// reads the MongoDB URL from the environment variable MONGODB_URL
// reads the MongoDB Database from the environment variable MONGODB_DATABASE
func InitMongoDB() (err error) {
	// TODO: logging?
	mDbSession, err := mongo.NewClientWithOptions(
		os.Getenv("MONGODB_URL"),
		clientopt.AppName("Cacophony"),
		clientopt.ConnectTimeout(30*time.Second),
		//clientopt.Dialer(),
		clientopt.ReadConcern(readconcern.Majority()),
		clientopt.WriteConcern(writeconcern.New(writeconcern.WMajority())),
	)
	if err != nil {
		return err
	}

	err = mDbSession.Connect(context.Background())
	if err != nil {
		return err
	}

	cache.SetMongo(mDbSession.Database(os.Getenv("MONGODB_DATABASE")))
	return nil
}
