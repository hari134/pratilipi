package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/messaging"
	"github.com/hari134/pratilipi/userservice/internal/dto"
	"github.com/hari134/pratilipi/userservice/middleware"
	"github.com/hari134/pratilipi/userservice/models"
	"github.com/hari134/pratilipi/userservice/producer"
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

// GetUserByIdHandler handles HTTP GET requests to retrieve a user by ID.
func (h *UserAPIHandler) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL parameters
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	// Convert userID to int64
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch the user from the database
	var user models.User
	ctx := context.Background()
	err = h.DB.NewSelect().Model(&user).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	// Return the user as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetUsersHandler handles HTTP GET requests to retrieve all users.
func (h *UserAPIHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch all users from the database
	var users []models.User
	ctx := context.Background()
	err := h.DB.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// Return the list of users as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
