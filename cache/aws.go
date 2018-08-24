package cache

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	awsSession      *session.Session
	awsSessionMutex sync.RWMutex
)

// SetAwsSession caches an aws session for future use
func SetAwsSession(s *session.Session) {
	awsSessionMutex.Lock()
	defer awsSessionMutex.Unlock()

	awsSession = s
}

// GetAwsSession returns a cached aws session
func GetAwsSession() *session.Session {
	awsSessionMutex.RLock()
	defer awsSessionMutex.RUnlock()

	return awsSession
}
