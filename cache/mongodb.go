package cache

import (
	"sync"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var (
	mongoClient      *mongo.Database
	mongoClientMutex sync.RWMutex
)

// SetMongo caches a MongoDB Client for future use
func SetMongo(c *mongo.Database) {
	mongoClientMutex.Lock()
	defer mongoClientMutex.Unlock()

	mongoClient = c
}

// GetMongo returns a cached MongoDB Client
func GetMongo() *mongo.Database {
	mongoClientMutex.RLock()
	defer mongoClientMutex.RUnlock()

	return mongoClient
}
