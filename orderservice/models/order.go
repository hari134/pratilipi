package models

import (
	"time"
	"github.com/uptrace/bun"
)



type Order struct {
    bun.BaseModel `bun:"table:orders"`       // Map struct to "orders" table

    OrderID    int64     `bun:"order_id,pk,autoincrement"`  // Primary key
    UserID     int64     `bun:"user_id,notnull"`            // Reference to the user who placed the order
    TotalPrice float64   `bun:"total_price,notnull"`        // Total price of the order
    Status     string    `bun:"status,notnull"`             // Order status: placed, shipped, completed, etc.
    PlacedAt   time.Time `bun:"placed_at,default:current_timestamp"`  // Timestamp when the order was placed
    UpdatedAt  time.Time `bun:"updated_at,default:current_timestamp"` // Timestamp for the last update
}


type OrderItem struct {
    bun.BaseModel `bun:"table:order_items"`     // Map struct to "order_items" table

    OrderItemID   int64   `bun:"order_item_id,pk,autoincrement"`      // Primary key
    OrderID       int64   `bun:"order_id,notnull"`                   // Reference to the order
    ProductID     int64   `bun:"product_id,notnull"`                 // Reference to the product
    Quantity      int     `bun:"quantity,notnull"`                   // Quantity of the product ordered
    PriceAtOrder  float64 `bun:"price_at_order,notnull"`             // Product price at the time of the order
}
