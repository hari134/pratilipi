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
}

// VerifyClaims checks if the token is valid by calling the userservice validate-token API.
func verifyClaims(ctx context.Context, role string) (ValidateTokenResponse, error) {
	// Extract the token from the context
	token, ok := ctx.Value("authtoken").(string)
	if !ok || token == "" {
		return ValidateTokenResponse{}, fmt.Errorf("unauthorized: token not provided")
	}

	// Call the validate-token API to check if the token is valid
	validateURL := "http://userservice:8080/validate-token"
	reqBody, _ := json.Marshal(map[string]string{
		"token": token,
	})

	resp, err := http.Post(validateURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		return ValidateTokenResponse{}, fmt.Errorf("unauthorized: invalid token")
	}
	defer resp.Body.Close()

	// Decode the response from the userservice
	var validateResp ValidateTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil || !validateResp.Valid {
		return ValidateTokenResponse{}, fmt.Errorf("unauthorized: invalid token")
	}
	if validateResp.Role == "admin" {
		return validateResp, nil
	}
	fmt.Println(validateResp.Role)
	if validateResp.Role != role {
		return ValidateTokenResponse{}, fmt.Errorf("unauthorized: role does not pemit this operation")
	}
	return validateResp, nil
}

// Users is the resolver for the users query.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	_, err := verifyClaims(ctx, "user")
	if err != nil {
		return nil, err
	}

	resp, err := http.Get("http://userservice:8080/users")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []struct {
		UserID int64  `json:"userID"`
		Name   string `json:"name"`
		Email  string `json:"email"`
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
	_, err := verifyClaims(ctx, "user")
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("http://userservice:8080/users/%s", id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user struct {
		UserID int64  `json:"userID"`
		Name   string `json:"name"`
		Email  string `json:"email"`
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
	_, err := verifyClaims(ctx, "user")
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
	_, err := verifyClaims(ctx, "user")
	if err != nil {
		return nil, err
	}

	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID")
	}
	fmt.Println(productID)
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

func (r *queryResolver) Orders(ctx context.Context) ([]*model.Order, error) {
	// Verify claims for user role
	_, err := verifyClaims(ctx, "user")
	if err != nil {
		return nil, err
	}

	// Fetch orders from the external service
	resp, err := http.Get("http://orderservice:8080/orders")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Structure to hold the API response
	var ordersWithItems []struct {
		OrderID    int64   `json:"OrderID"`
		UserID     int64   `json:"UserID"`
		TotalPrice float64 `json:"TotalPrice"`
		Status     string  `json:"Status"`
		PlacedAt   string  `json:"PlacedAt"`
		UpdatedAt  string  `json:"UpdatedAt"`
		OrderItems []struct {
			ProductID int64 `json:"ProductID"`
			Quantity  int   `json:"Quantity"`
		} `json:"OrderItems"`
	}

	// Decode the API response
	err = json.NewDecoder(resp.Body).Decode(&ordersWithItems)
	if err != nil {
		return nil, err
	}

	// Map the response to the GraphQL model
	var orders []*model.Order
	for _, apiOrder := range ordersWithItems {
		order := &model.Order{
			OrderID:    strconv.FormatInt(apiOrder.OrderID, 10), // Convert int64 to string
			UserID:     strconv.FormatInt(apiOrder.UserID, 10),  // Convert int64 to string
			TotalPrice: apiOrder.TotalPrice,
			Status:     apiOrder.Status,
			PlacedAt:   apiOrder.PlacedAt,
			Items:      []*model.OrderItem{}, // Initialize the items slice
		}

		// Map order items
		for _, apiItem := range apiOrder.OrderItems {
			gqlItem := &model.OrderItem{
				ProductID: strconv.FormatInt(apiItem.ProductID, 10), // Convert int64 to string
				Quantity:  apiItem.Quantity,
			}
			order.Items = append(order.Items, gqlItem)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *queryResolver) Order(ctx context.Context, id string) (*model.Order, error) {
	// Call the external orderservice API
	resp, err := http.Get(fmt.Sprintf("http://orderservice:8080/orders/%s", id))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %v", err)
	}
	defer resp.Body.Close()

	// Structure to hold the API response
	var orderWithItems struct {
		OrderID    int64   `json:"OrderID"`
		UserID     int64   `json:"UserID"`
		TotalPrice float64 `json:"TotalPrice"`
		Status     string  `json:"Status"`
		PlacedAt   string  `json:"PlacedAt"`
		UpdatedAt  string  `json:"UpdatedAt"`
		OrderItems []struct {
			ProductID    int64   `json:"ProductID"`
			Quantity     int     `json:"Quantity"`
			PriceAtOrder float64 `json:"PriceAtOrder"`
		} `json:"OrderItems"`
	}

	// Decode the API response
	err = json.NewDecoder(resp.Body).Decode(&orderWithItems)
	if err != nil {
		return nil, fmt.Errorf("failed to decode order response: %v", err)
	}

	// Map the API response to the GraphQL model
	order := &model.Order{
		OrderID:    strconv.FormatInt(orderWithItems.OrderID, 10), // Convert int64 to string for GraphQL ID type
		UserID:     strconv.FormatInt(orderWithItems.UserID, 10),  // Convert int64 to string for GraphQL ID type
		TotalPrice: orderWithItems.TotalPrice,
		Status:     orderWithItems.Status,
		PlacedAt:   orderWithItems.PlacedAt,
		Items:      []*model.OrderItem{}, // Initialize the slice
	}

	// Map the OrderItems to GraphQL model
	for _, item := range orderWithItems.OrderItems {
		gqlItem := &model.OrderItem{
			ProductID: strconv.FormatInt(item.ProductID, 10), // Convert int64 to string
			Quantity:  item.Quantity,
		}
		order.Items = append(order.Items, gqlItem)
	}

	return order, nil
}

// Mutation resolvers

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterInput) (*model.User, error) {
	// Marshal the input (name, phone_no, email, password, role) into JSON
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("could not marshal input: %v", err)
	}
	fmt.Println(input.Role)
	// Send the POST request to the user service for registering a new user
	resp, err := http.Post("http://userservice:8080/create-user", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to user service: %v", err)
	}
	defer resp.Body.Close()

	// Define a temporary struct to hold the API response with `UserID` as an int64
	var apiResponse struct {
		UserID       int64  `json:"UserID"`
		Name         string `json:"Name"`
		Email        string `json:"Email"`
		PhoneNo      string `json:"PhoneNo"`
		PasswordHash string `json:"PasswordHash"`
		Role         string `json:"Role"`
		CreatedAt    string `json:"CreatedAt"`
		UpdatedAt    string `json:"UpdatedAt"`
	}

	// Decode the response from the user service
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Prepare the GraphQL response object
	var user model.User
	user.UserID = strconv.FormatInt(apiResponse.UserID, 10) // Convert int64 to string
	user.Email = apiResponse.Email
	user.PhoneNo = apiResponse.PhoneNo
	user.Name = apiResponse.Name

	// Return the newly registered user
	return &user, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.ProductInput) (*model.Product, error) {
	_, err := verifyClaims(ctx, "admin")
	if err != nil {
		return nil, err
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

// PlaceOrder is the resolver for the placeOrder mutation.
func (r *mutationResolver) PlaceOrder(ctx context.Context, input model.OrderInput) (*model.Order, error) {
	claims, err := verifyClaims(ctx, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to get claims: %v", err)
	}

	input.UserID = strconv.FormatInt(claims.UserID, 10)
	fmt.Println(input.UserID)
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post("http://orderservice:8080/orders", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Order Service: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	var order struct {
		OrderID    int64              `json:"OrderID"`
		UserID     int64 `json:"UserID"`
		Items      []*model.OrderItem `json:"Items"`
		TotalPrice float64            `json:"TotalPrice"`
		Status     string             `json:"Status"`
		PlacedAt   string             `json:"PlacedAt"` // Should be a string representing a timestamp
	}
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response from Order Service: %v", err)
	}
	var gqlOrder = model.Order{
		OrderID: strconv.FormatInt(order.OrderID,10),
		UserID: strconv.FormatInt(order.UserID,10),
		Items : order.Items,
		TotalPrice: order.TotalPrice,
		PlacedAt: order.PlacedAt,
	}
	return &gqlOrder, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
