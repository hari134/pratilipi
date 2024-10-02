package dto

// ValidateTokenRequest represents the request body for token validation.
type ValidateTokenRequest struct {
	Token string `json:"token"`
}

// ValidateTokenResponse represents the response body for token validation.
type ValidateTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserID int64  `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
	Role   string `json:"role,omitempty"`
	Error  string `json:"error,omitempty"`
}

// LoginRequest represents the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response body for user login.
type LoginResponse struct {
	Token string `json:"token"`
}
