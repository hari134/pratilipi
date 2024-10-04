package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/messaging"
	"github.com/hari134/pratilipi/productservice/models"
	"github.com/hari134/pratilipi/productservice/producer"
)

type ProductAPIHandler struct {
	DB       *db.DB
	Producer *producer.ProducerManager
}

// CreateProductHandler handles the creation of a new product and emits a ProductCreated event.
func (h *ProductAPIHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	ctx := context.Background()
	_, err := h.DB.NewInsert().Model(&product).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	productIDStr := strconv.FormatInt(product.ProductID, 10)
	// Emit the ProductCreated event to Kafka
	event := &messaging.ProductCreated{
		ProductID:      productIDStr,
		Name:           product.Name,
		Price:          product.Price,
		InventoryCount: product.InventoryCount,
	}

	if err := h.Producer.EmitProductCreatedEvent(event); err != nil {
		http.Error(w, "Failed to emit ProductCreated event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// UpdateProductHandler updates a product's details and emits a ProductUpdated event.
func (h *ProductAPIHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["product_id"]

	var productUpdate models.Product

	if err := json.NewDecoder(r.Body).Decode(&productUpdate); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	product := &models.Product{}
	err := h.DB.NewSelect().Model(product).Where("product_id = ?", productID).Scan(ctx)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Update product fields
	product.Name = productUpdate.Name
	product.Description = productUpdate.Description
	product.Price = productUpdate.Price
	product.UpdatedAt = time.Now()

	_, err = h.DB.NewUpdate().Model(product).Where("product_id = ?", productID).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// DeleteProductHandler deletes a product and emits a ProductDeleted event.
func (h *ProductAPIHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["product_id"]

	ctx := context.Background()
	product := &models.Product{}
	err := h.DB.NewSelect().Model(product).Where("product_id = ?", productID).Scan(ctx)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	_, err = h.DB.NewDelete().Model(product).Where("product_id = ?", productID).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
}

// UpdateInventoryHandler updates the inventory of a product and emits an InventoryUpdated event.
func (h *ProductAPIHandler) UpdateInventoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["product_id"]

	var inventoryUpdate struct {
		InventoryCount int `json:"inventory_count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&inventoryUpdate); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	product := &models.Product{}
	err := h.DB.NewSelect().Model(product).Where("product_id = ?", productID).Scan(ctx)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	product.InventoryCount = inventoryUpdate.InventoryCount
	product.UpdatedAt = time.Now()

	_, err = h.DB.NewUpdate().Model(product).Where("product_id = ?", productID).Exec(ctx)
	if err != nil {
		http.Error(w, "Failed to update inventory", http.StatusInternalServerError)
		return
	}

	productIDStr := strconv.FormatInt(product.ProductID, 10)

	// Emit the InventoryUpdated event to Kafka
	event := &messaging.ProductInventoryUpdated{
		ProductID:      productIDStr,
		InventoryCount: product.InventoryCount,
	}

	if err := h.Producer.EmitInventoryUpdatedEvent(event); err != nil {
		http.Error(w, "Failed to emit InventoryUpdated event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
