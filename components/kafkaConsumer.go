package components

import (
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// InitKafkaConsumer initializes and caches a Kafka Consumer client
// reads the list of kafka brokers from the environment variable KAFKA_BROKER (delimited by ,)
func InitKafkaConsumer() error {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	saramaConsumer, err := sarama.NewConsumer(strings.Split(os.Getenv("KAFKA_BROKER"), ","), saramaConfig)
	if err != nil {
		return err
	}
	cache.SetKafkaConsumer(saramaConsumer)

	return nil
}
