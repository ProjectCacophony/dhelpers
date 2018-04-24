package cache

import (
	"sync"

	"github.com/globalsign/mgo"
)

var (
	mgoSession      *mgo.Database
	mgoSessionMutex sync.RWMutex
)

// SetMgo caches a MongoDB Database Session for future use
func SetMgo(s *mgo.Database) {
	mgoSessionMutex.Lock()
	defer mgoSessionMutex.Unlock()

	mgoSession = s
}

// GetMgo returns a cached MongoDB Database Session
func GetMgo() *mgo.Database {
	mgoSessionMutex.RLock()
	defer mgoSessionMutex.RUnlock()

	return mgoSession
}
