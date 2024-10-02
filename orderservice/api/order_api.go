package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/hari134/pratilipi/orderservice/models"
	"github.com/hari134/pratilipi/pkg/db"
)

// OrderHandler handles order-related API requests.
type OrderHandler struct {
	DB *db.DB  // Injected database dependency
}

// OrderRequest represents the payload for placing an order.
type OrderRequest struct {
	UserID string          `json:"user_id"`
	Items  []OrderItemData `json:"items"`
}

// OrderItemData represents an individual item in the order.
type OrderItemData struct {
	ProductID    string  `json:"product_id"`
	Quantity     int     `json:"quantity"`
	PriceAtOrder float64 `json:"price_at_order"`
}

// PlaceOrderHandler handles the HTTP POST request to place an order.
func (h *OrderHandler) PlaceOrderHandler(w http.ResponseWriter, r *http.Request) {
	var orderReq OrderRequest

	// Decode the incoming request body into OrderRequest struct
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate that the user exists in the users table
	ctx := context.Background()
	var user models.User
	err := h.DB.NewSelect().Model(&user).Where("user_id = ?", orderReq.UserID).Scan(ctx)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Validate that products exist and check stock
	for _, item := range orderReq.Items {
		var product models.Product
		err := h.DB.NewSelect().Model(&product).Where("product_id = ?", item.ProductID).Scan(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("Product %s not found", item.ProductID), http.StatusNotFound)
			return
		}

		// Check if enough stock is available
		if product.Stock < item.Quantity {
			http.Error(w, fmt.Sprintf("Insufficient stock for product %s. Available: %d, Requested: %d",
				product.ProductID, product.Stock, item.Quantity), http.StatusBadRequest)
			return
		}
	}

	// Create the order in the orders table
	order := &models.Order{
		UserID:    orderReq.UserID,
		TotalPrice: calculateTotalPrice(orderReq.Items), // Calculate total price from items
		Status:     "placed",
		PlacedAt:   time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = h.DB.NewInsert().Model(order).Exec(ctx)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Create order items in the order_items table
	for _, item := range orderReq.Items {
		orderItem := &models.OrderItem{
			OrderID:      order.OrderID,  // Reference the order just created
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			PriceAtOrder: item.PriceAtOrder,
			StockAtOrder: getStockForProduct(item.ProductID, h.DB), // Get stock at the time of order placement
		}

		_, err := h.DB.NewInsert().Model(orderItem).Exec(ctx)
		if err != nil {
			log.Printf("Failed to create order item: %v", err)
			http.Error(w, "Failed to create order items", http.StatusInternalServerError)
			return
		}
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// calculateTotalPrice calculates the total price of the order based on the items.
func calculateTotalPrice(items []OrderItemData) float64 {
	var totalPrice float64
	for _, item := range items {
		totalPrice += item.PriceAtOrder * float64(item.Quantity)
	}
	return totalPrice
}

// getStockForProduct retrieves the stock of the product from the database.
func getStockForProduct(productID string, db *db.DB) int {
	ctx := context.Background()
	var product models.Product
	err := db.NewSelect().Model(&product).Where("product_id = ?", productID).Scan(ctx)
	if err != nil {
		return 0
	}
	return product.Stock
}
