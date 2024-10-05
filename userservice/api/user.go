package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/messaging"
	"github.com/hari134/pratilipi/userservice/internal/dto"
	"github.com/hari134/pratilipi/userservice/middleware"
	"github.com/hari134/pratilipi/userservice/models"
	"github.com/hari134/pratilipi/userservice/producer"
	"golang.org/x/crypto/bcrypt"
)

// UserAPIHandler holds dependencies for the user API routes.
type UserAPIHandler struct {
	DB            *db.DB
	KafkaProducer *producer.ProducerManager
}

// UserRequest represents the incoming payload for creating a user, including the plain password.
type UserRequest struct {
	Name     string `json:"name"`
	PhoneNo  string `json:"phone_no"`
	Email    string `json:"email"`
	Password string `json:"password"`  // Plain password
	Role     string `json:"role"`
}

// CreateUserHandler handles HTTP POST requests to create a new user and stores a hashed password.
func (h *UserAPIHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userReq UserRequest

	// Decode the incoming request body into UserRequest struct
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the user's password
	hashedPassword, err := hashPassword(userReq.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Prepare the user model to be inserted into the database
	user := models.User{
		Name:         userReq.Name,
		PhoneNo:      userReq.PhoneNo,
		Email:        userReq.Email,
		PasswordHash: hashedPassword,
		Role:         userReq.Role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	// Insert the user into the database
	ctx := context.Background()
	_, err = h.DB.NewInsert().Model(&user).Exec(ctx)
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

	// Respond with the created user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// hashPassword hashes a plain password using bcrypt.
func hashPassword(password string) (string, error) {
	// Use bcrypt to generate a hashed password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
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
