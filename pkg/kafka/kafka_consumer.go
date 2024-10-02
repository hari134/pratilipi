package kafka

import (
    "context"
    "encoding/json"
    "log"
    "github.com/hari134/pratilipi/pkg/messaging"
    "github.com/segmentio/kafka-go"
)

// KafkaConsumer implements the Consumer interface for Kafka.
type KafkaConsumer struct {
    Reader *kafka.Reader
}

// NewKafkaConsumer creates a new Kafka consumer using the provided KafkaConfig.
func NewKafkaConsumer(config *KafkaConfig) messaging.Consumer {
    return &KafkaConsumer{
        Reader: kafka.NewReader(kafka.ReaderConfig{
            Brokers:  config.Brokers,       // Use brokers from config
            GroupID:  config.GroupID,       // Use group ID from config
            Topic:    config.Topic,         // Use topic from config
            MinBytes: 10e3,  // 10KB
            MaxBytes: 10e6,  // 10MB
        }),
    }
}

// Subscribe listens to messages from the specified Kafka topic and deserializes them into the eventStruct.
func (kc *KafkaConsumer) Subscribe(topic string, eventStruct interface{}, handler func(event interface{}) error) error {
    kc.Reader.SetOffset(kafka.FirstOffset)
    for {
        msg, err := kc.Reader.ReadMessage(context.Background())
        if err != nil {
            log.Printf("Error reading message from topic %s: %v", topic, err)
            continue
        }

        eventInstance := eventStruct
        if err := json.Unmarshal(msg.Value, eventInstance); err != nil {
            log.Printf("Failed to unmarshal message: %v", err)
            continue
        }

        if err := handler(eventInstance); err != nil {
            log.Printf("Error processing message: %v", err)
            continue
        }
    }
}

// Close gracefully closes the Kafka consumer connection.
func (kc *KafkaConsumer) Close() error {
    return kc.Reader.Close()
}
