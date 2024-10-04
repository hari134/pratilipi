package models

import (
    "time"
    "github.com/uptrace/bun"
)

type User struct {
    bun.BaseModel `bun:"table:users"`    // Map struct to "users" table

    UserID    int64     `bun:"user_id,pk"`                           // User ID (received from the User Service)
    Email     string    `bun:"email,notnull"`                        // Email of the user (received from the User Service)
    PhoneNo   string    `bun:"phone_no"`                             // Phone number of the user (optional, received from the User Service)
    CreatedAt time.Time `bun:"created_at,default:current_timestamp"` // Timestamp when the user was added
}
