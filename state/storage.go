package state

import (
	"sync"

	"time"

	"github.com/go-redis/redis"
	"github.com/json-iterator/go"
	"gitlab.com/Cacophony/dhelpers/cache"
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

func guildBannedUserIDsSetKey(guildID string) string {
	return "project-d:state:guild-" + guildID + ":banned-userids"
}

func guildBannedUserIDInitializedGuildIDsSetKey() string {
	return "project-d:state:banned-userids-initialized-guild-ids"
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

func messagesListKey(channelID string) string {
	return "project-d:state:channel-" + channelID + ":messages"
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

func addToStateSet(key string, items ...string) (err error) {
	interfaceItems := make([]interface{}, 0)
	for _, item := range items {
		interfaceItems = append(interfaceItems, item)
	}
	if len(interfaceItems) <= 0 {
		return
	}

	err = cache.GetRedisClient().SAdd(key, interfaceItems...).Err()
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

func addToStateList(key string, items ...string) (err error) {
	interfaceItems := make([]interface{}, 0)
	for _, item := range items {
		interfaceItems = append(interfaceItems, item)
	}
	if len(interfaceItems) <= 0 {
		return
	}

	err = cache.GetRedisClient().LPush(key, interfaceItems...).Err()
	return err
}

func trimStateList(key string, limit int64) (err error) {
	err = cache.GetRedisClient().LTrim(key, 0, limit).Err()
	return err
}

func readStateList(key string) (items []string, err error) {
	items, err = cache.GetRedisClient().LRange(key, 0, -1).Result()
	return items, err
}
