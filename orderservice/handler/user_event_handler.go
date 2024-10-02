package handler

import (
	"context"
	"log"
	"github.com/hari134/pratilipi/orderservice/models"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/messaging"
)

// UserEventHandler handles both UserRegistered and UserProfileUpdated events.
type UserEventHandler struct {
	DB *db.DB  // Injected database dependency
}

// NewUserEventHandler creates a new UserEventHandler instance with the provided database.
func NewUserEventHandler(dbInstance *db.DB) *UserEventHandler {
	return &UserEventHandler{DB: dbInstance}
}

// HandleUserRegistered processes the UserRegistered event received from Kafka.
func (h *UserEventHandler) HandleUserRegistered(event *messaging.UserRegistered) error {
	log.Printf("Processing UserRegistered event: %+v", event)

	// Insert the new user into the database
	ctx := context.Background()
	user := &models.User{
		UserID:  event.UserID,
		PhoneNo: event.PhoneNo,
		Email:   event.Email,
	}

	// Insert user into the DB using the injected Bun ORM instance
	_, err := h.DB.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}

	log.Printf("User %s inserted successfully", event.Email)
	return nil
}

// HandleUserProfileUpdated processes the UserProfileUpdated event and updates the user in the Order Service's DB.
func (h *UserEventHandler) HandleUserProfileUpdated(event *messaging.UserProfileUpdated) error {
	log.Printf("Processing UserProfileUpdated event: %+v", event)

	// Update the user in the users table in the Order Service's DB
	ctx := context.Background()
	user := &models.User{
		UserID:  event.UserID,
		Email:   event.Email,
		PhoneNo: event.PhoneNo,
	}

	_, err := h.DB.NewUpdate().Model(user).Where("user_id = ?", event.UserID).Exec(ctx)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return err
	}

	log.Printf("User %s updated successfully", event.Email)
	return nil
}
