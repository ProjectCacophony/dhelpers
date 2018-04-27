package cache

import (
	"sync"

	"github.com/minio/minio-go"
)

var (
	minioClient      *minio.Client
	minioClientMutex sync.RWMutex
)

// SetMinio caches a minio client for future use
func SetMinio(s *minio.Client) {
	minioClientMutex.Lock()
	defer minioClientMutex.Unlock()

	minioClient = s
}

// GetMinio returns a cached minio client
func GetMinio() *minio.Client {
	minioClientMutex.Lock()
	defer minioClientMutex.Unlock()

	return minioClient
}
