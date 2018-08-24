package cache

import (
	"sync"

	"github.com/Shopify/sarama"
)

var (
	kafkaConsumer      sarama.Consumer
	kafkaConsumerMutex sync.RWMutex
)

// SetKafkaConsumer caches a Kafka Consumer client for future use
func SetKafkaConsumer(s sarama.Consumer) {
	kafkaConsumerMutex.Lock()
	defer kafkaConsumerMutex.Unlock()

	kafkaConsumer = s
}

// GetKafkaConsumer returns a cached Kafka Consumer client
func GetKafkaConsumer() sarama.Consumer {
	kafkaConsumerMutex.RLock()
	defer kafkaConsumerMutex.RUnlock()

	return kafkaConsumer
}
