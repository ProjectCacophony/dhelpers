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
	if err != nil {
		return err
	}

	//fmt.Println("saved", key, "to state", "(size: "+humanize.Bytes(uint64(binary.Size(marshalled)))+")")
	return nil
}

func deleteStateObject(key string) error {
	err := cache.GetRedisClient().Del(key).Err()
	if err != nil {
		return err
	}

	//fmt.Println("deleted", key, "from state")
	return nil
}

func readStateObject(key string) (data []byte, err error) {
	data, err = cache.GetRedisClient().Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrStateNotFound
		}
		return nil, err
	}

	//fmt.Println("read", key, "from state", "(size: "+humanize.Bytes(uint64(binary.Size(data)))+")")
	return data, nil
}
