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
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

// ProductInventoryUpdated represents the event when product inventory is updated.
type ProductInventoryUpdated struct {
	ProductID string `json:"product_id"`
	Stock     int    `json:"stock"`
}

// OrderPlaced event is emitted when a user places an order.
type OrderPlaced struct {
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	ProductID   string    `json:"product_id"`
	Quantity    int       `json:"quantity"`
	TotalAmount float64   `json:"total_amount"`
	PlacedAt    time.Time `json:"placed_at"`
}

// OrderShipped event is emitted when an order is shipped.
type OrderShipped struct {
	OrderID   string    `json:"order_id"`
	ShippedAt time.Time `json:"shipped_at"`
}
