package models

import (
	"github.com/uptrace/bun"
)

// Product represents a product in the Order Service's local cache.
type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ProductID string  `bun:"product_id,pk"`         // Unique identifier for the product
	Name      string  `bun:"name"`                  // Name of the product
	Price     float64 `bun:"price"`                 // Product price at creation time
	Stock     int     `bun:"stock"`                 // Current stock level
	CreatedAt string  `bun:"created_at,default:current_timestamp"`
	UpdatedAt string  `bun:"updated_at,default:current_timestamp on update current_timestamp"`
}
