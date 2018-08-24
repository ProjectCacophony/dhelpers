package components

import (
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
	"github.com/wvanbergen/kazoo-go"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// InitKafkaConsumerGroup initializes and caches a Kafka Consumer Group client
// reads the list of zookeepers from the environment variable ZOOKEEPER_ADDRESSES (delimited by ,)
func InitKafkaConsumerGroup() error {
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true
	config.Zookeeper.Logger = cache.GetLogger()

	var zookeeperNodes []string
	zookeeperNodes, config.Zookeeper.Chroot = kazoo.ParseConnectionString(
		strings.Trim(os.Getenv("ZOOKEEPER_ADDRESSES"), ","),
	)

	kafkaConsumerGroup, err := consumergroup.JoinConsumerGroup(
		"cacophony", []string{"cacophony"}, zookeeperNodes, config,
	)
	if err != nil {
		return err
	}

	cache.SetKafkaConsumerGroup(kafkaConsumerGroup)

	return nil
}
