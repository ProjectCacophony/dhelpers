package dhelpers

import (
	"testing"

	"os"

	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gitlab.com/Cacophony/dhelpers/cache"
)

func init() {
	// init logger
	cache.SetLogger(logrus.NewEntry(logrus.New()))
	// init redis
	cache.SetRedisClient(redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB:       0,
	}))
}

func TestGetEventKey(t *testing.T) {
	v := GetEventKey(&discordgo.MessageCreate{
		Message: &discordgo.Message{
			ID:        "foo",
			ChannelID: "bar",
			GuildID:   "foo",
			Content:   "bar",
		},
	})
	if v != "cacophony:gateway:event-MESSAGE_CREATE-4bba92e91bad1523affa8369b63171fc" {
		t.Error("Expected cacophony:gateway:event-MESSAGE_CREATE-4bba92e91bad1523affa8369b63171fc, got ", v)
	}
}

func TestIsNewEvent(t *testing.T) {
	key1 := "cacophony:gateway:event-MESSAGE_CREATE-" + strconv.FormatInt(time.Now().Unix(), 10)
	key2 := "cacophony:gateway:event-MESSAGE_UPDATE-" + strconv.FormatInt(time.Now().Unix(), 10)
	v := IsNewEvent(cache.GetRedisClient(), "testing", key1)
	if !v {
		t.Error("Expected true, got ", v)
	}
	v = IsNewEvent(cache.GetRedisClient(), "testing", key1)
	if v {
		t.Error("Expected false, got ", v)
	}
	v = IsNewEvent(cache.GetRedisClient(), "testing", key2)
	if !v {
		t.Error("Expected true, got ", v)
	}
}
