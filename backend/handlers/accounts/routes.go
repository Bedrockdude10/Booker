// handlers/accounts/routes.go
package accounts

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
Modern router with clean API design - no legacy support
*/
func Routes(r chi.Router, collections map[string]*mongo.Collection) {
	service := NewService(collections)
	jwtService := NewJWTService()
	handler := NewHandler(service, jwtService)

	// Public authentication routes
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/register", handler.Register)
		r.Post("/refresh", handler.RefreshToken)
	})

	// Protected account management routes
	r.Route("/api/account", func(r chi.Router) {
		r.Use(handler.AuthMiddleware)

		r.Get("/", handler.GetAccount)                     // Get current user
		r.Put("/", handler.UpdateAccount)                  // Update current user
		r.Post("/change-password", handler.ChangePassword) // Change password
	})

	// Admin-only routes
	r.Route("/api/admin/accounts", func(r chi.Router) {
		r.Use(handler.AuthMiddleware)
		r.Use(handler.AdminMiddleware)

		r.Get("/", handler.ListAccounts)                 // List all accounts
		r.Get("/{id}", handler.GetAccount)               // Get account by ID
		r.Put("/{id}", handler.UpdateAccount)            // Update any account
		r.Delete("/{id}", handler.DeactivateAccount)     // Deactivate account
		r.Put("/{id}/activate", handler.ActivateAccount) // Reactivate account
		r.Put("/{id}/password", handler.ChangePassword)  // Admin password reset
	})
}
