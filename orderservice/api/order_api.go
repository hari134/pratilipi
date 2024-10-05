package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/hari134/pratilipi/orderservice/models"
	"github.com/hari134/pratilipi/orderservice/producer"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/messaging"
)

// OrderHandler handles order-related API requests.
type OrderHandler struct {
	DB       *db.DB
	Producer *producer.ProducerManager
}

// OrderRequest represents the payload for placing an order.
type OrderRequest struct {
	UserID string          `json:"user_id"`
	Items  []*OrderItemData `json:"items"`
}

// OrderItemData represents an individual item in the order request.
type OrderItemData struct {
	ProductID    int64   `json:"product_id"`
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

	// Validate that products exist and check stock
	for _, item := range orderReq.Items {
		var product models.Product
		err := h.DB.NewSelect().Model(&product).Where("product_id = ?", item.ProductID).Scan(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("Product %d not found", item.ProductID), http.StatusNotFound)
			return
		}

		// Check if enough stock is available
		if product.InventoryCount < item.Quantity {
			fmt.Printf("Insufficient stock for product %d. Available: %d, Requested: %d",
				product.ProductID, product.InventoryCount, item.Quantity)
			http.Error(w, fmt.Sprintf("Insufficient stock for product %d. Available: %d, Requested: %d",
				product.ProductID, product.InventoryCount, item.Quantity), http.StatusBadRequest)
			return
		}
	}
	userId, err := strconv.ParseInt(orderReq.UserID, 10, 64)
	if err != nil {
		panic(err)
	}
	// Create the order in the orders table
	order := &models.Order{
		UserID:     userId,
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
	var eventItems []messaging.OrderItem
	var orderItemsArr []models.OrderItem
	for _, item := range orderReq.Items {
		orderItem := &models.OrderItem{
			OrderID:      order.OrderID, // Reference the order just created
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			PriceAtOrder: item.PriceAtOrder,
		}

		_, err := h.DB.NewInsert().Model(orderItem).Exec(ctx)
		if err != nil {
			log.Printf("Failed to create order item: %v", err)
			http.Error(w, "Failed to create order items", http.StatusInternalServerError)
			return
		}
		orderItemsArr = append(orderItemsArr,*orderItem)
		// Collect items for the event
		eventItems = append(eventItems, messaging.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})

		// Update the product stock
		product := &models.Product{
			ProductID: item.ProductID,
		}
		_, err = h.DB.NewUpdate().
			Model(product).
			Set("inventory_count = inventory_count - ?", item.Quantity).
			Where("product_id = ?", item.ProductID).
			Exec(ctx)
		if err != nil {
			log.Printf("Failed to update product stock: %v", err)
			http.Error(w, "Failed to update product stock", http.StatusInternalServerError)
			return
		}
	}

	// Emit "Order Placed" event
	orderPlacedEvent := &messaging.OrderPlaced{
		OrderID: order.OrderID,
		UserID:  order.UserID,
		Items:   eventItems,
	}
	err = h.Producer.EmitOrderPlacedEvent(orderPlacedEvent)
	if err != nil {
		log.Printf("Failed to emit OrderPlaced event: %v", err)
		http.Error(w, "Failed to emit OrderPlaced event", http.StatusInternalServerError)
		return
	}

	// Return success response
	order.OrderItems = orderItemsArr
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetAllOrdersHandler handles the HTTP GET request to retrieve all orders with their items.
func (h *OrderHandler) GetAllOrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Fetch all orders first
	var orders []models.Order
	err := h.DB.NewSelect().Model(&orders).Order("placed_at DESC").Scan(ctx)
	if err != nil {
		log.Printf("Failed to retrieve orders: %v", err)
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}

	// Create a map to hold the orders and their items
	ordersWithItems := make([]struct {
		Order      models.Order
		OrderItems []models.OrderItem
	}, len(orders))

	// For each order, fetch its associated order items
	for i, order := range orders {
		ordersWithItems[i].Order = order

		var orderItems []models.OrderItem
		err := h.DB.NewSelect().Model(&orderItems).Where("order_id = ?", order.OrderID).Scan(ctx)
		if err != nil {
			log.Printf("Failed to retrieve items for order %d: %v", order.OrderID, err)
			http.Error(w, "Failed to retrieve order items", http.StatusInternalServerError)
			return
		}

		ordersWithItems[i].OrderItems = orderItems
	}

	// Return the list of orders with items
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordersWithItems)
}

// GetOrderByIDHandler handles the HTTP GET request to retrieve a specific order by its ID with its items.
func (h *OrderHandler) GetOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	// Fetch the order
	var order models.Order
	err := h.DB.NewSelect().
		Model(&order).
		Where("order_id = ?", orderID).
		Scan(ctx)
	if err != nil {
		log.Printf("Failed to retrieve order with ID %s: %v", orderID, err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Fetch items associated with the order
	var orderItems []models.OrderItem
	err = h.DB.NewSelect().
		Model(&orderItems).
		Where("order_id = ?", order.OrderID).
		Scan(ctx)
	if err != nil {
		log.Printf("Failed to retrieve items for order with ID %s: %v", orderID, err)
		http.Error(w, "Failed to retrieve order items", http.StatusInternalServerError)
		return
	}

	// Attach the items to the order
	order.OrderItems = orderItems

	// Return the order with its items
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// calculateTotalPrice calculates the total price of the order based on the items.
func calculateTotalPrice(items []*OrderItemData) float64 {
	var totalPrice float64
	for _, item := range items {
		totalPrice += item.PriceAtOrder * float64(item.Quantity)
	}
	return totalPrice
}
