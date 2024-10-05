package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hari134/pratilipi/graphqlgateway/graph/model"
)

// Query resolvers

// Users is the resolver for the users query.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	// Fetch the users from the external service
	resp, err := http.Get("http://userservice:8080/users")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Define a temporary struct to hold the data with userID as int64
	var users []struct {
		UserID       int64  `json:"userID"`
		Name         string `json:"name"`
		Email        string `json:"email"`
	}

	// Decode the response into the temporary struct
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, err
	}

	// Convert userID to string and map to GraphQL model
	var gqlUsers []*model.User
	for _, user := range users {
		gqlUser := &model.User{
			UserID: strconv.FormatInt(user.UserID, 10), // Convert userID from int64 to string
			Name:   user.Name,
			Email:  user.Email,
		}
		gqlUsers = append(gqlUsers, gqlUser)
	}

	// Return the list of users
	return gqlUsers, nil
}

// User is the resolver for the user query.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	// Fetch the user by ID from the external service
	resp, err := http.Get(fmt.Sprintf("http://userservice:8080/users/%s", id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Define a temporary struct to hold the data with userID as int64
	var user struct {
		UserID       int64  `json:"userID"`
		Name         string `json:"name"`
		Email        string `json:"email"`
	}

	// Decode the response into the temporary struct
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	// Convert userID to string and map to GraphQL model
	gqlUser := &model.User{
		UserID: strconv.FormatInt(user.UserID, 10), // Convert userID from int64 to string
		Name:   user.Name,
		Email:  user.Email,
	}

	// Return the user
	return gqlUser, nil
}

// Products is the resolver for the products query.
func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	// Fetch products from external service
	resp, err := http.Get("http://productservice:8080/products")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Define temporary struct to match the Product model with ProductID as int64
	var products []struct {
		ProductID      int64     `json:"productID"`
		Name           string    `json:"name"`
		Description    string    `json:"description"`
		Price          float64   `json:"price"`
		InventoryCount int       `json:"inventoryCount"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	// Convert to GraphQL model
	var gqlProducts []*model.Product
	for _, product := range products {
		gqlProduct := &model.Product{
			ProductID:      strconv.FormatInt(product.ProductID, 10), // Convert int64 to string
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			InventoryCount: product.InventoryCount,
			CreatedAt:      product.CreatedAt.Format(time.RFC3339),   // Format time as string
			UpdatedAt:      product.UpdatedAt.Format(time.RFC3339),   // Format time as string
		}
		gqlProducts = append(gqlProducts, gqlProduct)
	}

	return gqlProducts, nil
}

// Product is the resolver for the product query.
func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	// Convert the ID from string to int64 for querying the service
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID")
	}

	// Fetch the product from the external service by ID
	resp, err := http.Get(fmt.Sprintf("http://productservice:8080/products/%d", productID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Define a temporary struct to hold the data with productID as int64
	var product struct {
		ProductID      int64     `json:"productID"`
		Name           string    `json:"name"`
		Description    string    `json:"description"`
		Price          float64   `json:"price"`
		InventoryCount int       `json:"inventoryCount"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}

	// Decode the response into the temporary struct
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return nil, err
	}

	// Convert to GraphQL model
	gqlProduct := &model.Product{
		ProductID:      strconv.FormatInt(product.ProductID, 10), // Convert int64 to string
		Name:           product.Name,
		Description:    product.Description,
		Price:          product.Price,
		InventoryCount: product.InventoryCount,
		CreatedAt:      product.CreatedAt.Format(time.RFC3339),   // Format time as string
		UpdatedAt:      product.UpdatedAt.Format(time.RFC3339),   // Format time as string
	}

	return gqlProduct, nil
}

// Orders is the resolver for the orders query.
func (r *queryResolver) Orders(ctx context.Context) ([]*model.Order, error) {
	resp, err := http.Get("http://orderservice:8080/orders")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orders []*model.Order
	err = json.NewDecoder(resp.Body).Decode(&orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Order is the resolver for the order query.
func (r *queryResolver) Order(ctx context.Context, id string) (*model.Order, error) {
	resp, err := http.Get(fmt.Sprintf("http://orderservice:8080/orders/%s", id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var order model.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Mutation resolvers

// RegisterUser is the resolver for the registerUser mutation.
func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterInput) (*model.User, error) {
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://userservice:8080/register", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user model.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateProduct is the resolver for the createProduct mutation.
func (r *mutationResolver) CreateProduct(ctx context.Context, input model.ProductInput) (*model.Product, error) {
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://productservice:8080/products", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var product model.Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// PlaceOrder is the resolver for the placeOrder mutation.
func (r *mutationResolver) PlaceOrder(ctx context.Context, input model.OrderInput) (*model.Order, error) {
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://orderservice:8080/orders", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var order model.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}


// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
