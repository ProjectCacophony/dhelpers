package components

import (
	"os"

	"github.com/go-redis/redis"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// InitRedis initializes and caches the redis client
// reads the redis address from the environment variable REDIS_ADDRESS, example: 127.0.0.1:6379
func InitRedis() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB:       0,
	})
	cache.SetRedisClient(redisClient)
}
