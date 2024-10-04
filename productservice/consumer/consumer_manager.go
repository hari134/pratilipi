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

// StartOrderPlacedConsumer listens for "Order Placed" events and updates product inventory.
func (cm *ConsumerManager) StartOrderPlacedConsumer(topic string) {
    log.Printf("Starting consumer for topic: %s", topic)

    // Set up a handler function to process "Order Placed" events.
    handler := func(event interface{}) error {
        orderPlaced := event.(*messaging.OrderPlaced)
        return cm.handleOrderPlaced(orderPlaced)
    }

    // Subscribe to the "order-placed" topic and process incoming events.
    if err := cm.consumer.Subscribe(topic, &messaging.OrderPlaced{}, handler); err != nil {
        log.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
    }
}

// handleOrderPlaced processes the "Order Placed" event and updates the inventory for each product.
func (cm *ConsumerManager) handleOrderPlaced(event *messaging.OrderPlaced) error {
    log.Printf("Processing OrderPlaced event: %+v", event)

    ctx := context.Background()

    // Loop through each item in the order and update the product inventory.
    for _, item := range event.Items {
        product := &models.Product{}
        err := cm.DB.NewSelect().Model(product).Where("product_id = ?", item.ProductID).Scan(ctx)
        if err != nil {
            log.Printf("Failed to find product with ID %d: %v", item.ProductID, err)
            return err
        }

        // Check if there's enough inventory to fulfill the order.
        if product.InventoryCount < item.Quantity {
            log.Printf("Not enough inventory for product %d", item.ProductID)
            return nil // Optionally, you could return an error here to handle this case.
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
