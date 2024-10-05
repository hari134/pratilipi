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

type ClaimsCtxKey struct{}

// Define the structure to hold the token validation response
type ValidateTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Error  string `json:"error"`
}
// Utility function to extract JWT claims from the context
func getClaims(ctx context.Context) (ValidateTokenResponse, error) {
	claims, ok := ctx.Value(ClaimsCtxKey{}).(ValidateTokenResponse)
	if !ok {
		return ValidateTokenResponse{}, fmt.Errorf("unauthorized: invalid token or no token provided")
	}
	return claims, nil
}

// Users is the resolver for the users query.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	claims, err := getClaims(ctx)
	if err != nil {
		return nil, err
	}

	// Only admins can fetch all users
	if claims.Role != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can view users")
	}

	resp, err := http.Get("http://userservice:8080/users")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []struct {
		UserID       int64  `json:"userID"`
		Name         string `json:"name"`
		Email        string `json:"email"`
	}

	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, err
	}

	var gqlUsers []*model.User
	for _, user := range users {
		gqlUser := &model.User{
			UserID: strconv.FormatInt(user.UserID, 10),
			Name:   user.Name,
			Email:  user.Email,
		}
		gqlUsers = append(gqlUsers, gqlUser)
	}

	return gqlUsers, nil
}

// User is the resolver for the user query.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	claims, err := getClaims(ctx)
	if err != nil {
		return nil, err
	}

	if strconv.FormatInt(claims.UserID, 10) != id && claims.Role != "admin" {
		return nil, fmt.Errorf("unauthorized: you can only view your own profile")
	}

	resp, err := http.Get(fmt.Sprintf("http://userservice:8080/users/%s", id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user struct {
		UserID       int64  `json:"userID"`
		Name         string `json:"name"`
		Email        string `json:"email"`
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	gqlUser := &model.User{
		UserID: strconv.FormatInt(user.UserID, 10),
		Name:   user.Name,
		Email:  user.Email,
	}

	return gqlUser, nil
}

// Products is the resolver for the products query.
func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	_, err := getClaims(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get("http://productservice:8080/products")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []struct {
		ProductID      int64     `json:"productID"`
		Name           string    `json:"name"`
		Description    string    `json:"description"`
		Price          float64   `json:"price"`
		InventoryCount int       `json:"inventoryCount"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}

	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	var gqlProducts []*model.Product
	for _, product := range products {
		gqlProduct := &model.Product{
			ProductID:      strconv.FormatInt(product.ProductID, 10),
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			InventoryCount: product.InventoryCount,
			CreatedAt:      product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      product.UpdatedAt.Format(time.RFC3339),
		}
		gqlProducts = append(gqlProducts, gqlProduct)
	}

	return gqlProducts, nil
}

// Product is the resolver for the product query.
func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	_, err := getClaims(ctx)
	if err != nil {
		return nil, err
	}

	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID")
	}

	resp, err := http.Get(fmt.Sprintf("http://productservice:8080/products/%d", productID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var product struct {
		ProductID      int64     `json:"productID"`
		Name           string    `json:"name"`
		Description    string    `json:"description"`
		Price          float64   `json:"price"`
		InventoryCount int       `json:"inventoryCount"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}

	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return nil, err
	}

	gqlProduct := &model.Product{
		ProductID:      strconv.FormatInt(product.ProductID, 10),
		Name:           product.Name,
		Description:    product.Description,
		Price:          product.Price,
		InventoryCount: product.InventoryCount,
		CreatedAt:      product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      product.UpdatedAt.Format(time.RFC3339),
	}

	return gqlProduct, nil
}

// Orders is the resolver for the orders query.
func (r *queryResolver) Orders(ctx context.Context) ([]*model.Order, error) {
	claims, err := getClaims(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get("http://orderservice:8080/orders")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ordersWithItems []struct {
		Order      model.Order      `json:"order"`
		OrderItems []*model.OrderItem `json:"order_items"`
	}

	err = json.NewDecoder(resp.Body).Decode(&ordersWithItems)
	if err != nil {
		return nil, err
	}

	var orders []*model.Order
	for _, orderWithItems := range ordersWithItems {
		order := orderWithItems.Order
		order.Items = orderWithItems.OrderItems
		if claims.Role == "admin" || strconv.FormatInt(claims.UserID, 10) == order.UserID {
			orders = append(orders, &order)
		}
	}

	return orders, nil
}

// Order is the resolver for the order query.
func (r *queryResolver) Order(ctx context.Context, id string) (*model.Order, error) {
	claims, err := getClaims(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("http://orderservice:8080/orders/%s", id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orderWithItems struct {
		Order      model.Order      `json:"order"`
		OrderItems []*model.OrderItem `json:"order_items"`
	}

	err = json.NewDecoder(resp.Body).Decode(&orderWithItems)
	if err != nil {
		return nil, err
	}

	order := orderWithItems.Order
	order.Items = orderWithItems.OrderItems

	if claims.Role == "admin" || strconv.FormatInt(claims.UserID, 10) == order.UserID {
		return &order, nil
	}

	return nil, fmt.Errorf("unauthorized: you can only view your own orders")
}

// Mutation resolvers

// RegisterUser is the resolver for the registerUser mutation.
func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterInput) (*model.User, error) {
	// Marshal the input (name, phone_no, email, password, role) into JSON
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("could not marshal input: %v", err)
	}

	// Send the POST request to the user service for registering a new user
	resp, err := http.Post("http://userservice:8080/register", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to user service: %v", err)
	}
	defer resp.Body.Close()

	// Decode the response from the user service to get the newly registered user
	var user model.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Return the newly registered user
	return &user, nil
}


func (r *mutationResolver) CreateProduct(ctx context.Context, input model.ProductInput) (*model.Product, error) {
	claims, err := getClaims(ctx)
	if err != nil {
		return nil, err
	}

	// Only admins are allowed to create products
	if claims.Role != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can create products")
	}

	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	// Send the request to create the product
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
	claims, err := getClaims(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get claims: %v", err)
	}

	input.UserID = strconv.FormatInt(claims.UserID, 10)

	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post("http://orderservice:8080/orders", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Order Service: %v", err)
	}
	defer resp.Body.Close()

	var order model.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response from Order Service: %v", err)
	}

	return &order, nil
}



// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }