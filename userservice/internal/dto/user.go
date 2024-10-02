package dto


type UpdateUserRequest struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    Name   string `json:"name"`
}