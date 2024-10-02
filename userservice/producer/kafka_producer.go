package producer

import (
    "encoding/json"
    "log"
    "github.com/hari134/pratilipi/pkg/messaging"
    "github.com/hari134/pratilipi/pkg/kafka"
)

// KafkaProducerManager manages Kafka producers.
type KafkaProducerManager struct {
    producer *kafka.KafkaProducer
}

// NewKafkaProducerManager creates a new instance of KafkaProducerManager.
func NewKafkaProducerManager(producer *kafka.KafkaProducer) *KafkaProducerManager {
    return &KafkaProducerManager{producer: producer}
}

// EmitUserRegisteredEvent emits a UserRegistered event to Kafka.
func (kpm *KafkaProducerManager) EmitUserRegisteredEvent(event *messaging.UserRegistered) error {
    eventBytes, err := json.Marshal(event)
    if err != nil {
        return err
    }

    log.Printf("Emitting UserRegistered event: %s", eventBytes)
    return kpm.producer.Emit("user-registered", eventBytes)
}

// EmitUserProfileUpdatedEvent emits a UserProfileUpdated event to Kafka.
func (kpm *KafkaProducerManager) EmitUserProfileUpdatedEvent(event *messaging.UserProfileUpdated) error {
    eventBytes, err := json.Marshal(event)
    if err != nil {
        return err
    }

    log.Printf("Emitting UserProfileUpdated event: %s", eventBytes)
    return kpm.producer.Emit("user-profile-updated", eventBytes)
}

