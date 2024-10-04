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
	cm.subscribeToTopics([]string{userRegisteredTopic, productCreatedTopic})
}

// subscribeToTopics handles events from multiple topics with a unified event handler.
func (cm *ConsumerManager) subscribeToTopics(topics []string) {
	for _, topic := range topics {
		// Subscribe each topic with a common handler
		go func(topic string) {
			handler := func(event interface{}) error {
				return cm.handleEvent(event, topic)
			}

			if err := cm.consumer.Subscribe(topic, topic, handler); err != nil {
				log.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
			}
		}(topic)
	}
}

// handleEvent processes different event types using a switch case.
func (cm *ConsumerManager) handleEvent(event interface{}, topic string) error {
	switch e := event.(type) {
	case *messaging.UserRegistered:
		// Handle User Registered event
		return cm.handleUserRegistered(e)
	case *messaging.ProductCreated:
		// Handle Product Created event
		return cm.handleProductCreated(e)
	default:
		log.Printf("Received unknown event type on topic %s: %+v", topic, event)
		// Optionally, handle unknown event types if necessary
		return nil
	}
}

// handleUserRegistered processes the "User Registered" event and updates the users table.
func (cm *ConsumerManager) handleUserRegistered(event *messaging.UserRegistered) error {
	log.Printf("Processing UserRegistered event: %+v", event)

	ctx := context.Background()

	userIdInt, err := strconv.ParseInt(event.UserID, 10, 64)
	if err != nil {
		return err
	}
	user := &models.User{
		UserID:  userIdInt,
		Email:   event.Email,
		PhoneNo: event.PhoneNo,
	}

	_, err = cm.DB.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}

	log.Printf("User %s inserted successfully", event.Email)
	return nil
}

// handleProductCreated processes the "Product Created" event and updates the products table.
func (cm *ConsumerManager) handleProductCreated(event *messaging.ProductCreated) error {
	log.Printf("Processing ProductCreated event: %+v", event)

	ctx := context.Background()
	productIdInt, err := strconv.ParseInt(event.ProductID, 10, 64)
	if err != nil {
		return err
	}
	product := &models.Product{
		ProductID:      productIdInt,
		Price:          event.Price,
		InventoryCount: event.InventoryCount,
	}

	_, err = cm.DB.NewInsert().Model(product).Exec(ctx)
	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		return err
	}

	log.Printf("Product %s inserted successfully", event.Name)
	return nil
}
