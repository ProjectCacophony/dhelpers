package components

import (
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// InitKafkaProducer initializes and caches a Kafka Producer client
// reads the list of kafka brokers from the environment variable KAFKA_BROKER (delimited by ,)
func InitKafkaProducer() error {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	saramaConfig.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	saramaProducer, err := sarama.NewSyncProducer(strings.Split(os.Getenv("KAFKA_BROKER"), ","), saramaConfig)
	if err != nil {
		return err
	}
	cache.SetKafkaProducer(saramaProducer)

	return nil
}
