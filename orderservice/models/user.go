package models

import (
    "github.com/uptrace/bun"
)

// User represents a simplified user in the order service for validation purposes.
type User struct {
    bun.BaseModel `bun:"table:users,alias:u"`  // Bun ORM table mapping

    UserID   string `bun:"user_id,pk"`         // User ID (primary key)
    Email    string `bun:"email,notnull"`      // User's email
    PhoneNo  string `bun:"phone_no,notnull"`   // User's phone number
}
