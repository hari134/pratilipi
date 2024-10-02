package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/hari134/pratilipi/userservice/internal/jwtutil"
)

type contextKey string

const UserIDKey = contextKey("userID")

// TokenValidationMiddleware validates JWT token and extracts the user ID from it.
func TokenValidationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenStr := r.Header.Get("Authorization")
        if tokenStr == "" {
            http.Error(w, "Authorization header missing", http.StatusUnauthorized)
            return
        }

        // Strip 'Bearer' from the token string if it exists
        tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

        // Validate the token and extract claims
        claims, err := jwtutil.ParseJWTToken(tokenStr)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add the userID to the context
        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
