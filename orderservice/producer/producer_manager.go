package producer

import (
    "encoding/json"
    "log"
    "github.com/hari134/pratilipi/pkg/messaging"
)

// ProducerManager manages Kafka producers for the Order Service.
type ProducerManager struct {
    producer messaging.Producer
}

// NewProducerManager creates a new instance of ProducerManager.
func NewProducerManager(producer messaging.Producer) *ProducerManager {
    return &ProducerManager{producer: producer}
}

// EmitOrderPlacedEvent emits an OrderPlaced event to Kafka.
func (pm *ProducerManager) EmitOrderPlacedEvent(event *messaging.OrderPlaced) error {
    eventBytes, err := json.Marshal(event)
    if err != nil {
        return err
    }

    log.Printf("Emitting OrderPlaced event: %s", eventBytes)
    return pm.producer.Emit("order-placed", eventBytes)
}
