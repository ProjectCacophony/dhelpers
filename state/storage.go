package state

import (
	"sync"

	"time"

	"github.com/go-redis/redis"
	"github.com/json-iterator/go"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

var stateLock sync.Mutex
var stateExpire = time.Duration(0)

func guildUserIDsSetKey(guildID string) string {
	return "project-d:state:user-ids:" + guildID
}

func allUserIDsSetKey() string {
	return "project-d:state:user-ids"
}

func allGuildIDsSetKey() string {
	return "project-d:state:guild-ids"
}

func allChannelIDsSetKey() string {
	return "project-d:state:channel-ids"
}

func guildBotIDsSetKey(guildID string) string {
	return "project-d:state:guild-" + guildID + ":bot-ids"
}

func userKey(userID string) string {
	return "project-d:state:user-" + userID
}
func memberKey(guildID, userID string) string {
	return "project-d:state:guild-" + guildID + ":member-" + userID
}
func guildKey(guildID string) string {
	return "project-d:state:guild-" + guildID
}
func channelKey(channelID string) string {
	return "project-d:state:channel-" + channelID
}

func updateStateObject(key string, object interface{}) error {
	marshalled, err := jsoniter.Marshal(object)
	if err != nil {
		return err
	}

	err = cache.GetRedisClient().Set(key, marshalled, stateExpire).Err()
	return err
}

func deleteStateObject(key string) error {
	err := cache.GetRedisClient().Del(key).Err()
	return err
}

func readStateObject(key string) (data []byte, err error) {
	data, err = cache.GetRedisClient().Get(key).Bytes()

	if err == redis.Nil {
		return nil, ErrStateNotFound
	}

	return data, err
}

func addToStateSet(key, item string) (err error) {
	err = cache.GetRedisClient().SAdd(key, item).Err()
	return err
}

func removeFromStateSet(key, item string) (err error) {
	err = cache.GetRedisClient().SRem(key, item).Err()
	return err
}

func readStateSet(key string) (items []string, err error) {
	items, err = cache.GetRedisClient().SMembers(key).Result()
	return items, err
}
