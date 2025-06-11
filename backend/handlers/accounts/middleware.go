// handlers/accounts/middleware.go
package accounts

import (
	"context"
	"net/http"
	"strings"

	"github.com/Bedrockdude10/Booker/backend/utils"
)

// AuthMiddleware validates JWT tokens and sets user context
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.HandleError(w, utils.ValidationError("Authorization header required"))
			return
		}

		// Extract token from "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.HandleError(w, utils.ValidationError("Invalid authorization header format"))
			return
		}

		// Validate token
		claims, err := h.jwtService.ValidateToken(tokenParts[1])
		if err != nil {
			utils.HandleError(w, utils.ValidationError("Invalid or expired token"))
			return
		}

		// Add user claims to request context
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware validates JWT tokens if present but doesn't require them
func (h *Handler) OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// Extract token from "Bearer <token>" format
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
				// Validate token
				if claims, err := h.jwtService.ValidateToken(tokenParts[1]); err == nil {
					// Add user claims to request context
					ctx := context.WithValue(r.Context(), "user", claims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
		}

		// Continue without authentication
		next.ServeHTTP(w, r)
	})
}

// RoleMiddleware checks if user has required role
func (h *Handler) RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context (should be set by AuthMiddleware)
			claims, ok := r.Context().Value("user").(*Claims)
			if !ok {
				utils.HandleError(w, utils.ValidationError("User not found in context"))
				return
			}

			// Check role
			if claims.Role != requiredRole {
				utils.HandleError(w, utils.ValidationError("Insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AdminMiddleware is a convenience wrapper for admin-only routes
func (h *Handler) AdminMiddleware(next http.Handler) http.Handler {
	return h.RoleMiddleware("admin")(next)
}
