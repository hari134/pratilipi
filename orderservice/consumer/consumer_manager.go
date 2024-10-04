package consumer

import (
	"context"
	"log"
	"strconv"

	"github.com/hari134/pratilipi/orderservice/models"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/messaging"
)

// ConsumerManager listens for Kafka events and processes them for the Order Service.
type ConsumerManager struct {
	consumer messaging.Consumer
	DB       *db.DB // Injected database dependency
}

// NewConsumerManager creates a new instance of ConsumerManager.
func NewConsumerManager(consumer messaging.Consumer, dbInstance *db.DB) *ConsumerManager {
	return &ConsumerManager{
		consumer: consumer,
		DB:       dbInstance,
	}
}

// StartConsumers subscribes to the topics and processes different types of events.
func (cm *ConsumerManager) StartConsumers(userRegisteredTopic string, productCreatedTopic string) {
	// Subscribe to both topics with a unified handler
	handlers := map[string]func(event interface{}) error{
		userRegisteredTopic: cm.handleUserRegisteredEvent,
		productCreatedTopic: cm.handleProductCreatedEvent,
	}

	err := cm.consumer.Subscribe(handlers)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %v", err)
	}
}

// handleUserRegisteredEvent handles events from the "User Registered" topic.
func (cm *ConsumerManager) handleUserRegisteredEvent(event interface{}) error {
	log.Printf("Processing UserRegistered event: %+v", event)

	userRegistered, ok := event.(*messaging.UserRegistered)
	if !ok {
		log.Printf("Unexpected event type for UserRegistered event")
		return nil
	}

	ctx := context.Background()
	userIdInt, err := strconv.ParseInt(userRegistered.UserID, 10, 64)
	if err != nil {
		return err
	}
	user := &models.User{
		UserID:  userIdInt,
		Email:   userRegistered.Email,
		PhoneNo: userRegistered.PhoneNo,
	}

	_, err = cm.DB.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}

	log.Printf("User %s inserted successfully", userRegistered.Email)
	return nil
}

// handleProductCreatedEvent handles events from the "Product Created" topic.
func (cm *ConsumerManager) handleProductCreatedEvent(event interface{}) error {
	log.Printf("Processing ProductCreated event: %+v", event)

	productCreated, ok := event.(*messaging.ProductCreated)
	if !ok {
		log.Printf("Unexpected event type for ProductCreated event")
		return nil
	}

	ctx := context.Background()
	productIdInt, err := strconv.ParseInt(productCreated.ProductID, 10, 64)
	if err != nil {
		return err
	}
	product := &models.Product{
		ProductID:      productIdInt,
		Price:          productCreated.Price,
		InventoryCount: productCreated.InventoryCount,
	}

	_, err = cm.DB.NewInsert().Model(product).Exec(ctx)
	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		return err
	}

	log.Printf("Product %s inserted successfully", productCreated.Name)
	return nil
}
