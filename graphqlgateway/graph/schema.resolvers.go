package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/hari134/pratilipi/graphqlgateway/graph/model"
)

// Query resolvers

// Users is the resolver for the users query.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	resp, err := http.Get("http://userservice:8080/users")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []*model.User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// User is the resolver for the user query.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	resp, err := http.Get(fmt.Sprintf("http://userservice:8080/users/%s", id))
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

// Products is the resolver for the products query.
func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	resp, err := http.Get("http://productservice:8082/products")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []*model.Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Product is the resolver for the product query.
func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	resp, err := http.Get(fmt.Sprintf("http://productservice:8082/products/%s", id))
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

// Orders is the resolver for the orders query.
func (r *queryResolver) Orders(ctx context.Context) ([]*model.Order, error) {
	resp, err := http.Get("http://orderservice:8083/orders")
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
	resp, err := http.Get(fmt.Sprintf("http://orderservice:8083/orders/%s", id))
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

	resp, err := http.Post("http://userservice:8081/register", "application/json", bytes.NewBuffer(reqBody))
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

	resp, err := http.Post("http://productservice:8082/products", "application/json", bytes.NewBuffer(reqBody))
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

	resp, err := http.Post("http://orderservice:8083/orders", "application/json", bytes.NewBuffer(reqBody))
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
