package consumer

import (
    "context"
    "log"
    "github.com/hari134/pratilipi/pkg/db"
    "github.com/hari134/pratilipi/pkg/messaging"
    "github.com/hari134/pratilipi/productservice/models"
)

// ConsumerManager listens for events from Kafka and processes them.
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

// StartConsumers subscribes to the topics and processes different types of events.
func (cm *ConsumerManager) StartConsumers(orderPlacedTopic string) {
    // Subscribe to the "order-placed" topic with a unified handler.
    handlers := map[string]func(event interface{}) error{
        orderPlacedTopic: cm.handleOrderPlacedEvent,
    }

    err := cm.consumer.Subscribe(handlers)
    if err != nil {
        log.Fatalf("Failed to subscribe to topics: %v", err)
    }
}

// handleOrderPlacedEvent processes the "Order Placed" event and updates the inventory for each product.
func (cm *ConsumerManager) handleOrderPlacedEvent(event interface{}) error {
    log.Printf("Processing OrderPlaced event: %+v", event)

    orderPlaced, ok := event.(*messaging.OrderPlaced)
    if !ok {
        log.Printf("Unexpected event type for OrderPlaced event")
        return nil
    }

    ctx := context.Background()

    // Loop through each item in the order and update the product inventory.
    for _, item := range orderPlaced.Items {
        product := &models.Product{}
        err := cm.DB.NewSelect().Model(product).Where("product_id = ?", item.ProductID).Scan(ctx)
        if err != nil {
            log.Printf("Failed to find product with ID %d: %v", item.ProductID, err)
            return err
        }

        // Check if there's enough inventory to fulfill the order.
        if product.InventoryCount < item.Quantity {
            log.Printf("Not enough inventory for product %d", item.ProductID)
            return nil // Optionally, return an error here to handle this case.
        }

        // Deduct the quantity from the product's inventory.
        product.InventoryCount -= item.Quantity

        _, err = cm.DB.NewUpdate().Model(product).Where("product_id = ?", item.ProductID).Exec(ctx)
        if err != nil {
            log.Printf("Failed to update inventory for product %d: %v", item.ProductID, err)
            return err
        }

        log.Printf("Updated inventory for product %d: new inventory count is %d", item.ProductID, product.InventoryCount)
    }

    return nil
}
