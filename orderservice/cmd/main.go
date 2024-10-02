package main

import (
	"log"
	"time"

	"github.com/hari134/pratilipi/orderservice/consumer"
	"github.com/hari134/pratilipi/orderservice/handler"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/kafka"
	"github.com/hari134/pratilipi/pkg/messaging"
)

func main() {
	// Initialize the database
	dbConfig := db.Config{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "password",
		DBName:          "orderdb",
		SSLMode:         "disable",
		MaxIdleConns:    10,
		MaxOpenConns:    20,
		ConnMaxLifetime: 30 * time.Minute,
	}
	dbInstance := db.InitDB(dbConfig)
	defer db.CloseDB(dbInstance)

	// Create a new handler for UserRegistered with the injected database
	userEventHandler := handler.NewUserEventHandler(dbInstance)
	productEventHandler := handler.NewProductEventHandler(dbInstance)

	// Kafka configuration
	kafkaConfig := &kafka.KafkaConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "order-service-group",
	}

	// Define the handler configurations for each topic
	handlerConfigs := []consumer.HandlerConfig{
		{
			Topic: "user-registered",
			Handler: func(event interface{}) error {
				// Perform type assertion from interface{} to *messaging.UserRegistered
				if e, ok := event.(*messaging.UserRegistered); ok {
					return userEventHandler.HandleUserRegistered(e)
				}
				return nil // or handle the type assertion failure
			},
			Event: &messaging.UserRegistered{},
		},
		{
			Topic: "user-profile-updated",
			Handler: func(event interface{}) error {
				if e, ok := event.(*messaging.UserProfileUpdated); ok {
					return userEventHandler.HandleUserProfileUpdated(e)
				}
				return nil
			},
			Event: &messaging.UserProfileUpdated{},
		},
		{
			Topic: "product.created",
			Handler: func(event interface{}) error {
				if e, ok := event.(*messaging.ProductCreated); ok {
					return productEventHandler.HandleProductCreated(e)
				}
				return nil
			},
			Event: &messaging.ProductCreated{},
		},
		{
			Topic: "product.inventory.updated",
			Handler: func(event interface{}) error {
				if e, ok := event.(*messaging.ProductInventoryUpdated); ok {
					return productEventHandler.HandleProductInventoryUpdated(e)
				}
				return nil
			},
			Event: &messaging.ProductInventoryUpdated{},
		},
	}

	// Start Kafka consumers and pass the list of handler configurations
	consumer.StartKafkaConsumers(kafkaConfig, handlerConfigs)

	// Application running
	log.Println("Kafka consumers started. Application is running...")
}
