package messaging

import "time"

// UserRegistered event is emitted when a new user is registered.
type UserRegistered struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	PhoneNo string `json:"phone_no"`
}

// UserProfileUpdated event is emitted when a user updates their profile information.
type UserProfileUpdated struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	PhoneNo   string    `json:"phone_no,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductCreated struct {
	ProductID      string  `json:"product_id"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	InventoryCount int     `json:"inventory_count"`
}

// ProductInventoryUpdated represents the event when product inventory is updated.
type ProductInventoryUpdated struct {
	ProductID      string `json:"product_id"`
	InventoryCount int    `json:"inventory_count"`
}

// OrderPlaced represents the event for an order that has been placed.
type OrderPlaced struct {
    OrderID int64 `json:"order_id"`
    UserID  int64 `json:"user_id"`
    Items   []OrderItem `json:"items"`
}

// OrderItem represents a single item in an order.
type OrderItem struct {
    ProductID int64 `json:"product_id"`
    Quantity  int   `json:"quantity"`
}

// OrderShipped event is emitted when an order is shipped.
type OrderShipped struct {
	OrderID   string    `json:"order_id"`
	ShippedAt time.Time `json:"shipped_at"`
}
