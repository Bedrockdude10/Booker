package accounts

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Router maps endpoints to handlers for account-related operations
*/
func Routes(r chi.Router, collections map[string]*mongo.Collection) {
	service := NewService(collections)
	handler := &Handler{service: service}

	// Mount account routes under /api/accounts
	r.Route("/api/accounts", func(r chi.Router) {
		// Account management
		r.Post("/", handler.CreateAccount)           // Register new account
		r.Get("/{id}", handler.GetAccount)           // Get account by ID
		r.Put("/{id}", handler.UpdateAccount)        // Update account
		r.Delete("/{id}", handler.DeactivateAccount) // Soft delete account

		// Authentication routes
		r.Post("/login", handler.Login)                         // Login
		r.Post("/password/reset", handler.RequestPasswordReset) // Request password reset
		r.Put("/password/{id}", handler.UpdatePassword)         // Update password

		// Admin routes (you might want to add middleware for admin-only access)
		r.Get("/", handler.ListAccounts) // List all accounts (admin only)

		// Profile routes (for authenticated users to manage their own account)
		r.Get("/profile/{email}", handler.GetAccountByEmail)     // Get account by email
		r.Put("/profile/{id}/activate", handler.ActivateAccount) // Reactivate account
	})
}

// Health check specifically for accounts service
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	writeJSON(w, map[string]interface{}{
		"status":  "healthy",
		"service": "accounts",
	})
}
