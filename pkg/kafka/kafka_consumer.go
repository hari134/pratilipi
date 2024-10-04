package kafka

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/hari134/pratilipi/pkg/messaging"
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

func (kc *KafkaConsumer) Subscribe(topic string, eventName string, handler func(event interface{}) error) error {
	log.Printf("Subscribing to topic: %s", topic)

	for {
		msg, err := kc.Reader.FetchMessage(context.Background())
		if err != nil {
			log.Printf("Failed to fetch message from topic %s: %v", topic, err)
			return err
		}

		log.Printf("Message received from topic %s: %s", topic, string(msg.Value))

		// Check if the eventName exists in the TypeRegistry
		// eventType, exists := kc.TypeRegistry[eventName]
		// if !exists {
		// 	return errors.New("event type not registered")
		// }

		// Dynamically create a new instance of the registered event type
		var eventInstance messaging.UserRegistered
		// Base64 decode the message value (msg.Value)
		if err != nil {
			log.Printf("Failed to decode base64 string: %v", err)
			return err
		}
		// Unmarshal the message into the dynamically created event struct
		err = json.Unmarshal(msg.Value, eventInstance)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		// Call the event handler with the unmarshaled event
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
