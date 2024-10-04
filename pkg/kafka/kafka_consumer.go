package kafka

import (
    "github.com/segmentio/kafka-go"
)

// KafkaConsumer implements the Consumer interface for Kafka.
type KafkaConsumer struct {
    Reader *kafka.Reader
}

// NewKafkaConsumer creates a new Kafka consumer using the provided KafkaConfig.
func NewKafkaConsumer(config *KafkaConfig) *KafkaConsumer {
    readerConfig := kafka.ReaderConfig{
        Brokers:  config.Brokers,  // Use brokers from config
        GroupID:  config.GroupID,  // Use group ID from config
        MinBytes: 10e3,            // 10KB
        MaxBytes: 10e6,            // 10MB
    }

    // If multiple topics are set, use GroupTopics
    if len(config.GroupTopics) > 0 {
        readerConfig.GroupTopics = config.GroupTopics
    } else {
        readerConfig.Topic = config.Topic
    }

    return &KafkaConsumer{
        Reader: kafka.NewReader(readerConfig),
    }
}
