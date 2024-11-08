package models

import (
    "github.com/uptrace/bun"
)

type Product struct {
    bun.BaseModel `bun:"table:products"`  // Map struct to "products" table

    ProductID      int64     `bun:"product_id,pk"`                    // Product ID (received from the Product Service)
    Price          float64   `bun:"price,notnull"`                    // Product price
    InventoryCount int       `bun:"inventory_count,notnull"`          // Inventory count for the product
}
