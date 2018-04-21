package cache

import (
	"sync"

	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
	redisMutext sync.RWMutex
)

// SetRedisClient caches a redis client for future use
func SetRedisClient(s *redis.Client) {
	redisMutext.Lock()
	redisClient = s
	redisMutext.Unlock()
}

// GetRedisClient returns a cached redis client
func GetRedisClient() *redis.Client {
	redisMutext.RLock()
	defer redisMutext.RUnlock()

	return redisClient
}
