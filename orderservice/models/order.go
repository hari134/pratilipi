package models

import (
	"time"
	"github.com/uptrace/bun"
)

// Order represents an order in the system.
type Order struct {
	bun.BaseModel `bun:"table:orders,alias:o"`

	OrderID    int       `bun:"order_id,pk,autoincrement"`   // Primary key
	UserID     string    `bun:"user_id"`                    // User ID from User Service
	TotalPrice float64   `bun:"total_price"`                // Total price of the order
	Status     string    `bun:"status"`                     // Order status
	PlacedAt   time.Time `bun:"placed_at"`                  // Timestamp when the order was placed
	UpdatedAt  time.Time `bun:"updated_at"`                 // Timestamp for the last update
}

// OrderItem represents an individual item in an order.
type OrderItem struct {
	bun.BaseModel `bun:"table:order_items,alias:oi"`

	OrderItemID  int     `bun:"order_item_id,pk,autoincrement"` // Primary key
	OrderID      int     `bun:"order_id"`                      // Reference to orders table
	ProductID    string  `bun:"product_id"`                    // Product ID from Product Service
	Quantity     int     `bun:"quantity"`                      // Quantity of the product ordered
	PriceAtOrder float64 `bun:"price_at_order"`                // Product price at the time of the order
	StockAtOrder int     `bun:"stock_at_order"`                // Stock level at the time of the order (optional)
}
