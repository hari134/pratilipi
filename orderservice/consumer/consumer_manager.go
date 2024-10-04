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
    DB       *db.DB  // Injected database dependency
}

// NewConsumerManager creates a new instance of ConsumerManager.
func NewConsumerManager(consumer messaging.Consumer, dbInstance *db.DB) *ConsumerManager {
    return &ConsumerManager{
        consumer: consumer,
        DB:       dbInstance,
    }
}

// StartUserRegisteredConsumer listens for the "User Registered" event and updates the users table.
func (cm *ConsumerManager) StartUserRegisteredConsumer(topic string) {
    log.Printf("Starting consumer for topic: %s", topic)

    handler := func(event interface{}) error {
        userRegistered := event.(*messaging.UserRegistered)
        return cm.handleUserRegistered(userRegistered)
    }

    if err := cm.consumer.Subscribe(topic, &messaging.UserRegistered{}, handler); err != nil {
        log.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
    }
}

// handleUserRegistered processes the "User Registered" event and updates the users table.
func (cm *ConsumerManager) handleUserRegistered(event *messaging.UserRegistered) error {
    log.Printf("Processing UserRegistered event: %+v", event)

    ctx := context.Background()

    userIdInt , err := strconv.ParseInt(event.UserID,10,64)
    if err != nil{
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

// StartProductCreatedConsumer listens for the "Product Created" event and updates the products table.
func (cm *ConsumerManager) StartProductCreatedConsumer(topic string) {
    log.Printf("Starting consumer for topic: %s", topic)

    handler := func(event interface{}) error {
        productCreated := event.(*messaging.ProductCreated)
        return cm.handleProductCreated(productCreated)
    }

    if err := cm.consumer.Subscribe(topic, &messaging.ProductCreated{}, handler); err != nil {
        log.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
    }
}

// handleProductCreated processes the "Product Created" event and updates the products table.
func (cm *ConsumerManager) handleProductCreated(event *messaging.ProductCreated) error {
    log.Printf("Processing ProductCreated event: %+v", event)

    ctx := context.Background()
    productIdInt , err := strconv.ParseInt(event.ProductID,10,64)
    if err != nil{
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
