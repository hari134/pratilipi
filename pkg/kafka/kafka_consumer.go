package kafka

import (
	"context"
	"log"
	"reflect"

	"github.com/hari134/pratilipi/pkg/serde"
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

// Subscribe subscribes to a topic and processes messages using the provided handler.
func (kc *KafkaConsumer) Subscribe(topic string, eventStruct interface{}, handler func(event interface{}) error) error {
    log.Printf("Subscribing to topic: %s", topic)

    for {
        msg, err := kc.Reader.FetchMessage(context.Background())
        if err != nil {
            log.Printf("Failed to fetch message from topic %s: %v", topic, err)
            return err
        }
        log.Println("using reflect")
        log.Printf("Message received from topic %s: %s", topic, string(msg.Value))
        newEventStruct := reflect.New(reflect.TypeOf(eventStruct)).Interface()
        // Base64 decode and unmarshal the message into the event struct
        err = serde.Base64ToStruct(string(msg.Value), newEventStruct)
        if err != nil {
            log.Printf("Failed to decode and unmarshal message: %v", err)
            return err
        }

        // Call the event handler with the unmarshaled event
        err = handler(newEventStruct)
        if err != nil {
            log.Printf("Handler failed for topic %s: %v", topic, err)
            return err
        }

        // Commit the message after processing
        if err := kc.Reader.CommitMessages(context.Background(), msg); err != nil {
            log.Printf("Failed to commit message: %v", err)
            return err
        }
    }
}

// Close closes the Kafka consumer.
func (kc *KafkaConsumer) Close() error {
    return kc.Reader.Close()
}
