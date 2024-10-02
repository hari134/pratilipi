package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/userservice/internal/dto"
	"github.com/hari134/pratilipi/userservice/internal/jwtutil"
	"github.com/hari134/pratilipi/userservice/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthAPIHandler struct {
	DB *db.DB
}

func (h *AuthAPIHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	var user models.User
	err := h.DB.NewSelect().Model(&user).Where("email = ?", loginReq.Email).Scan(ctx)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := jwtutil.GenerateJWTToken(user.UserID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.LoginResponse{Token: token})
}

func (h *AuthAPIHandler) ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.ValidateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		http.Error(w, "Invalid request payload or missing token", http.StatusBadRequest)
		return
	}

	claims, err := jwtutil.ParseJWTToken(req.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(dto.ValidateTokenResponse{
			Valid: false,
			Error: "Invalid token",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ValidateTokenResponse{
		Valid:  true,
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	})
}
