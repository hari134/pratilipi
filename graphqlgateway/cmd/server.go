package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/hari134/pratilipi/graphqlgateway/graph"
)

const defaultPort = "8080"

// Define the key used for storing claims in context
type ClaimsCtxKey struct{}

// Define the structure to hold the token validation response
type ValidateTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// Add the JWT authentication middleware to the server
	http.Handle("/query", AuthMiddleware(srv))

	// GraphQL playground endpoint
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// AuthMiddleware is a middleware to extract the Authorization header and add it to the request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		// Set the token in the context, even if it's empty (as nil)
		ctx := context.WithValue(r.Context(), "authtoken", authHeader)
		r = r.WithContext(ctx)

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
