package producer

import (
    "encoding/json"
    "log"

    "github.com/hari134/pratilipi/pkg/messaging"
)

// ProducerManager manages any producer that implements the messaging.Producer interface.
type ProducerManager struct {
    producer messaging.Producer
}

// NewProducerManager creates a new instance of ProducerManager.
func NewProducerManager(producer messaging.Producer) *ProducerManager {
    return &ProducerManager{producer: producer}
}

// EmitProductCreatedEvent emits a ProductCreated event using the provided producer.
func (pm *ProducerManager) EmitProductCreatedEvent(event *messaging.ProductCreated) error {
    eventBytes, err := json.Marshal(event)
    if err != nil {
        return err
    }

    log.Printf("Emitting ProductCreated event: %s", eventBytes)
    return pm.producer.Emit("product-created", eventBytes)
}

// EmitInventoryUpdatedEvent emits an InventoryUpdated event using the provided producer.
func (pm *ProducerManager) EmitInventoryUpdatedEvent(event *messaging.ProductInventoryUpdated) error {
    eventBytes, err := json.Marshal(event)
    if err != nil {
        return err
    }

    log.Printf("Emitting InventoryUpdated event: %s", eventBytes)
    return pm.producer.Emit("inventory-updated", eventBytes)
}
