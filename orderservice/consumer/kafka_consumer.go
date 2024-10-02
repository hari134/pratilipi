package consumer

import (
	"log"
	"github.com/hari134/pratilipi/pkg/kafka"
)

// HandlerConfig holds a Kafka topic and its corresponding handler.
type HandlerConfig struct {
	Topic   string
	Handler func(event interface{}) error
	Event   interface{}  // The event structure to be used for the specific topic
}

// StartKafkaConsumers starts Kafka consumers for different topics and handlers.
func StartKafkaConsumers(kafkaConfig *kafka.KafkaConfig, handlerConfigs []HandlerConfig) {
	for _, hc := range handlerConfigs {
		// Start Kafka consumer for each topic and its corresponding handler
		go startConsumer(kafkaConfig, hc)
	}
}

// startConsumer listens to a given Kafka topic and processes messages using the provided handler.
func startConsumer(config *kafka.KafkaConfig, handlerConfig HandlerConfig) {
	consumer := kafka.NewKafkaConsumer(config)
	defer consumer.Close()

	eventHandler := func(event interface{}) error {
		// Call the handler for the event
		return handlerConfig.Handler(event)
	}

	log.Printf("Starting consumer for topic: %s", handlerConfig.Topic)
	if err := consumer.Subscribe(handlerConfig.Topic, handlerConfig.Event, eventHandler); err != nil {
		log.Fatalf("Failed to subscribe to topic %s: %v", handlerConfig.Topic, err)
	}
}
