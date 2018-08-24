package cache

import (
	"sync"

	"github.com/wvanbergen/kafka/consumergroup"
)

var (
	kafkaConsumerGroup      *consumergroup.ConsumerGroup
	kafkaConsumerGroupMutex sync.RWMutex
)

// SetKafkaConsumerGroup caches a Kafka Consumer Group clieent for future use
func SetKafkaConsumerGroup(s *consumergroup.ConsumerGroup) {
	kafkaConsumerGroupMutex.Lock()
	defer kafkaConsumerGroupMutex.Unlock()

	kafkaConsumerGroup = s
}

// GetKafkaConsumerGroup returns a cached Kafka Consumer Group client
func GetKafkaConsumerGroup() *consumergroup.ConsumerGroup {
	kafkaConsumerGroupMutex.RLock()
	defer kafkaConsumerGroupMutex.RUnlock()

	return kafkaConsumerGroup
}
