package cache

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logger      *logrus.Entry
	loggerMutex sync.RWMutex
)

func SetLogger(s *logrus.Entry) {
	loggerMutex.Lock()
	logger = s
	loggerMutex.Unlock()
}

func GetLogger() *logrus.Entry {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()

	if logger == nil {
		panic(errors.New("Tried to get logger before logger#SetLogger() was called"))
	}

	return logger
}
