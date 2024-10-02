package models

import (
	"time"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	UserID       int64     `bun:"user_id,pk,autoincrement"`  // Primary key
	Name         string    `bun:"name,notnull"`              // User's name
	PhoneNo      string    `bun:"phoneNo,notnull"`           // User's phone number
	Email        string    `bun:"email,unique,notnull"`      // Unique email address
	PasswordHash string    `bun:"password_hash,notnull"`     // Hashed password for authentication
	Role         string    `bun:"role,default:'user'"`       // Role: 'admin' or 'user', defaults to 'user'
	CreatedAt    time.Time `bun:"created_at,nullzero,default:current_timestamp"` // User registration timestamp
	UpdatedAt    time.Time `bun:"updated_at,nullzero,default:current_timestamp"` // Timestamp for last update
}
