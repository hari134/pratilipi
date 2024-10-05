package models

import "time"

// Product represents a product in the Product Service.
type Product struct {
    ProductID      int64     `bun:"product_id,pk,autoincrement"`  // Primary key
    Name           string    `bun:"name,notnull"`                 // Product name
    Description    string    `bun:"description,nullzero"`         // Product description
    Price          float64   `bun:"price,notnull"`                // Product price
    InventoryCount int       `bun:"inventory_count,notnull" json:"inventorycount"`      // Available inventory
    CreatedAt      time.Time `bun:"created_at,nullzero,default:current_timestamp"` // Timestamp when the product was created
    UpdatedAt      time.Time `bun:"updated_at,nullzero,default:current_timestamp"` // Timestamp for last update
}
