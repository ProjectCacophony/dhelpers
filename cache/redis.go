package cache

import (
	"sync"

	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
	redisMutext sync.RWMutex
)

func SetRedisClient(s *redis.Client) {
	redisMutext.Lock()
	redisClient = s
	redisMutext.Unlock()
}

func GetRedisClient() *redis.Client {
	redisMutext.RLock()
	defer redisMutext.RUnlock()

	return redisClient
}
