package handler

import (
	"context"
	"log"
	"github.com/hari134/pratilipi/orderservice/models"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/messaging"
)

// ProductEventHandler handles product-related events.
type ProductEventHandler struct {
	DB *db.DB  // Injected database dependency
}

func NewProductEventHandler(dbInstance *db.DB) *ProductEventHandler{
	return &ProductEventHandler{DB: dbInstance}
}
// HandleProductCreated processes the ProductCreated event received from Kafka.
func (h *ProductEventHandler) HandleProductCreated(event *messaging.ProductCreated) error {
	log.Printf("Processing ProductCreated event: %+v", event)

	// Insert the new product into the products table
	ctx := context.Background()
	product := &models.Product{
		ProductID: event.ProductID,
		Name:      event.Name,
		Price:     event.Price,
		Stock:     event.Stock,
	}

	_, err := h.DB.NewInsert().Model(product).Exec(ctx)
	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		return err
	}

	log.Printf("Product %s inserted successfully", event.Name)
	return nil
}

// HandleProductInventoryUpdated processes the ProductInventoryUpdated event received from Kafka.
func (h *ProductEventHandler) HandleProductInventoryUpdated(event *messaging.ProductInventoryUpdated) error {
	log.Printf("Processing ProductInventoryUpdated event: %+v", event)

	// Update the stock of the product in the products table
	ctx := context.Background()
	product := &models.Product{
		ProductID: event.ProductID,
		Stock:     event.Stock,
	}

	_, err := h.DB.NewUpdate().Model(product).Column("stock").Where("product_id = ?", event.ProductID).Exec(ctx)
	if err != nil {
		log.Printf("Failed to update product stock: %v", err)
		return err
	}

	log.Printf("Product stock updated successfully for product ID %s", event.ProductID)
	return nil
}
