package cache

import (
	"sync"

	"github.com/Shopify/sarama"
)

var (
	kafkaProducer      sarama.SyncProducer
	kafkaProducerMutex sync.RWMutex
)

// SetKafkaProducer caches a Kafka Producer client for future use
func SetKafkaProducer(s sarama.SyncProducer) {
	kafkaProducerMutex.Lock()
	defer kafkaProducerMutex.Unlock()

	kafkaProducer = s
}

// GetKafkaProducer returns a cached Kafka Producer client
func GetKafkaProducer() sarama.SyncProducer {
	kafkaProducerMutex.RLock()
	defer kafkaProducerMutex.RUnlock()

	return kafkaProducer
}
