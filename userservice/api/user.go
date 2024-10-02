package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/messaging"
	"github.com/hari134/pratilipi/userservice/models"
	"github.com/hari134/pratilipi/userservice/internal/dto"
	"github.com/hari134/pratilipi/userservice/producer"
	"github.com/hari134/pratilipi/userservice/middleware"
)

// UserAPIHandler holds dependencies for the user API routes.
type UserAPIHandler struct {
	DB            *db.DB
	KafkaProducer *producer.ProducerManager
}

// CreateUserHandler handles HTTP POST requests to create a new user and emits an event.
func (h *UserAPIHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	ctx := context.Background()
	_, err := h.DB.NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Emit the UserRegistered event to Kafka
	userIDStr := strconv.FormatInt(user.UserID, 10)
	event := &messaging.UserRegistered{
		UserID:  userIDStr,
		Email:   user.Email,
		PhoneNo: user.PhoneNo,
	}
	if err := h.KafkaProducer.EmitUserRegisteredEvent(event); err != nil {
		http.Error(w, "Failed to emit event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserAPIHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
    var updateReq dto.UpdateUserRequest

    // Decode the incoming request body
    if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Extract the userID from the JWT token via the context
    userIDFromToken := r.Context().Value(middleware.UserIDKey).(int64)

    // Ensure that the user is updating their own profile
    if userIDFromToken != updateReq.UserID {
        http.Error(w, "You can only update your own profile", http.StatusForbidden)
        return
    }

    // Proceed with updating the user's profile in the database
    ctx := context.Background()
    _, err := h.DB.NewUpdate().
        Model(&models.User{Email: updateReq.Email, Name: updateReq.Name}).
        Where("user_id = ?", updateReq.UserID).
        Exec(ctx)

    if err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "user updated"})
}