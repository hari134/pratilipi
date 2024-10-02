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

// EmitUserRegisteredEvent emits a UserRegistered event using the provided producer.
func (pm *ProducerManager) EmitUserRegisteredEvent(event *messaging.UserRegistered) error {
    eventBytes, err := json.Marshal(event)
    if err != nil {
        return err
    }

    log.Printf("Emitting UserRegistered event: %s", eventBytes)
    return pm.producer.Emit("user-registered", eventBytes)
}

// EmitUserProfileUpdatedEvent emits a UserProfileUpdated event using the provided producer.
func (pm *ProducerManager) EmitUserProfileUpdatedEvent(event *messaging.UserProfileUpdated) error {
    eventBytes, err := json.Marshal(event)
    if err != nil {
        return err
    }

    log.Printf("Emitting UserProfileUpdated event: %s", eventBytes)
    return pm.producer.Emit("user-profile-updated", eventBytes)
}
