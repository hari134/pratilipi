package kafka

import (
    "context"
    "encoding/json"
    "log"
    "github.com/hari134/pratilipi/pkg/messaging"
    "github.com/segmentio/kafka-go"
)

// KafkaProducer implements the Producer interface for Kafka.
type KafkaProducer struct {
    Writer *kafka.Writer
}

// NewKafkaProducer creates a new Kafka producer using the provided KafkaConfig.
func NewKafkaProducer(config *KafkaConfig) messaging.Producer {
    return &KafkaProducer{
        Writer: kafka.NewWriter(kafka.WriterConfig{
            Brokers: config.Brokers,        // Use brokers from config
            Balancer: &kafka.LeastBytes{},
        }),
    }
}

// Emit serializes the event struct to JSON and sends it to the specified Kafka topic.
func (kp *KafkaProducer) Emit(topic string, event interface{}) error {
    // Convert event struct to JSON
    eventBytes, err := json.Marshal(event)
    if err != nil {
        return err
    }

    // Send the JSON message to Kafka
    err = kp.Writer.WriteMessages(context.Background(), kafka.Message{
        Topic: topic,
        Value: eventBytes,
    })
    if err != nil {
        log.Printf("Failed to emit event to topic %s: %v", topic, err)
        return err
    }
    log.Printf("Message emitted to topic %s", topic)
    return nil
}

// Close gracefully closes the Kafka producer connection.
func (kp *KafkaProducer) Close() error {
    return kp.Writer.Close()
}
