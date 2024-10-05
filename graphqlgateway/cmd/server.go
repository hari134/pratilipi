package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

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
	Error  string `json:"error"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// Add the JWT authentication middleware to the server
	http.Handle("/query", jwtAuthMiddleware(srv))

	// GraphQL playground endpoint
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// jwtAuthMiddleware sends the token to userservice to validate the token
func jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// Token format: "Bearer <token>"
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Call the ValidateToken API in userservice
			validateURL := "http://userservice:8080/validate-token"
			reqBody, _ := json.Marshal(map[string]string{
				"token": token,
			})

			resp, err := http.Post(validateURL, "application/json", bytes.NewBuffer(reqBody))
			if err != nil || resp.StatusCode != http.StatusOK {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			defer resp.Body.Close()

			var validateResp ValidateTokenResponse
			if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil || !validateResp.Valid {
				// If token is invalid, return unauthorized error
				http.Error(w, validateResp.Error, http.StatusUnauthorized)
				return
			}

			// If token is valid, add the claims to the context
			ctx := context.WithValue(r.Context(), ClaimsCtxKey{}, validateResp)
			r = r.WithContext(ctx)
		} else {
			// If no token is provided, return unauthorized error
			http.Error(w, "Authorization token not provided", http.StatusUnauthorized)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
