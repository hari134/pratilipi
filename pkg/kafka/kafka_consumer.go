package kafka

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/segmentio/kafka-go"
)

// KafkaConsumer implements the Consumer interface for Kafka.
type KafkaConsumer struct {
	Reader       *kafka.Reader
	TypeRegistry map[string]reflect.Type
}

// NewKafkaConsumer creates a new Kafka consumer using the provided KafkaConfig.
func NewKafkaConsumer(config *KafkaConfig) *KafkaConsumer {
	readerConfig := kafka.ReaderConfig{
		Brokers:  config.Brokers, // Use brokers from config
		GroupID:  config.GroupID, // Use group ID from config
		MinBytes: 10e3,           // 10KB
		MaxBytes: 10e6,           // 10MB
	}

	// If multiple topics are set, use GroupTopics
	if len(config.GroupTopics) > 0 {
		readerConfig.GroupTopics = config.GroupTopics
	} else {
		readerConfig.Topic = config.Topic
	}

	return &KafkaConsumer{
		Reader:       kafka.NewReader(readerConfig),
		TypeRegistry: make(map[string]reflect.Type),
	}
}

func (kc *KafkaConsumer) RegisterType(eventName string, eventType interface{}) {
	kc.TypeRegistry[eventName] = reflect.TypeOf(eventType).Elem()
}

func (kc *KafkaConsumer) Subscribe(handlers map[string]func(event interface{}) error) error {
    log.Printf("Subscribing to multiple topics: %v", kc.Reader.Config().GroupTopics)

    for {
        msg, err := kc.Reader.FetchMessage(context.Background())
        if err != nil {
            log.Printf("Failed to fetch message: %v", err)
            return err
        }

        topic := msg.Topic // Determine the topic the message came from
        log.Printf("Message received from topic %s: %s", topic, string(msg.Value))

        // Step 1: Unmarshal the JSON string to extract the Base64 string
        var base64Str string
        err = json.Unmarshal(msg.Value, &base64Str)
        if err != nil {
            log.Printf("Failed to unmarshal JSON: %v", err)
            return err
        }

        // Step 2: Base64 decode the extracted string
        decodedValue, err := base64.StdEncoding.DecodeString(base64Str)
        if err != nil {
            log.Printf("Failed to decode base64 string: %v", err)
            return err
        }

        // Step 3: Determine the handler based on the topic or event type
        handler, exists := handlers[topic] // Get the handler based on the topic
        if !exists {
            log.Printf("No handler found for topic: %s", topic)
            return fmt.Errorf("no handler registered for topic: %s", topic)
        }

        // Step 4: Get the event type from the registered types
        eventType, exists := kc.TypeRegistry[topic]
        if !exists {
            return fmt.Errorf("event type not registered for topic: %s", topic)
        }
        eventInstance := reflect.New(eventType).Interface()

        // Step 5: Unmarshal the decoded JSON into the event struct
        err = json.Unmarshal(decodedValue, eventInstance)
        if err != nil {
            log.Printf("Failed to unmarshal decoded JSON: %v", err)
            return err
        }

        // Step 6: Call the appropriate handler with the unmarshaled event
        err = handler(eventInstance)
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
