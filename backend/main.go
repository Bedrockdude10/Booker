package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bedrockdude10/Booker/backend/handlers/artists"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title           Artist Recommendation API
// @version         1.0
// @description     API for music artist recommendations based on genre, location, and user preferences
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/Bedrockdude10/Booker
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	run(os.Stderr, os.Args[1:])
}

func run(stderr io.Writer, args []string) {
	// Parse command line flags
	cmd := flag.NewFlagSet("", flag.ExitOnError)
	verboseFlag := cmd.Bool("v", false, "Enable verbose logging")
	logLevelFlag := cmd.String("log-level", slog.LevelInfo.String(), "Log level (debug, info, warn, error)")
	if err := cmd.Parse(args); err != nil {
		fmt.Fprint(stderr, err)
		os.Exit(1)
	}

	// Set up structured logging
	logger := newLogger(*logLevelFlag, *verboseFlag, stderr)
	slog.SetDefault(logger)

	ctx := context.Background()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	// Get MongoDB connection string
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		// Default to localhost if not specified
		mongoURI = "mongodb://localhost:27017"
		slog.Info("Using default MongoDB URI", "uri", mongoURI)
	}

	// Get server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
		slog.Info("Using default port", "port", port)
	}

	// Set up MongoDB client
	slog.Info("Connecting to MongoDB...")
	client, err := connectToMongoDB(ctx, mongoURI)
	if err != nil {
		fatal(ctx, "Failed to connect to MongoDB", err)
	}
	defer func() {
		slog.Info("Disconnecting from MongoDB...")
		if err := client.Disconnect(ctx); err != nil {
			slog.Error("Failed to disconnect from MongoDB", "error", err)
		}
	}()

	// Set up database and collections
	db := client.Database("booker")
	collections := map[string]*mongo.Collection{
		"artists": db.Collection("artists"),
		// "userPreferences": db.Collection("userPreferences"),
	}

	// Set up Chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // You might want to restrict this in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not readily apparent
	}))

	// Mount artist routes
	artists.Routes(r, collections)

	// Add a simple health check route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine so it doesn't block
	go func() {
		slog.Info("Starting server", "port", port)
		slog.Info("Swagger UI available", "url", fmt.Sprintf("http://localhost:%s/swagger/index.html", port))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fatal(ctx, "Server failed", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	sig := <-quit
	slog.Info("Shutting down server", "signal", sig.String())

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		fatal(ctx, "Server forced to shutdown", err)
	}

	slog.Info("Server exited gracefully")
}

// newLogger creates a structured logger with the specified log level
func newLogger(logLevel string, verbose bool, stderr io.Writer) *slog.Logger {
	if verbose {
		logLevel = "debug"
	}

	level := slog.LevelInfo
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	return slog.New(slog.NewJSONHandler(stderr, &slog.HandlerOptions{
		AddSource: level == slog.LevelDebug,
		Level:     level,
	}))
}

// fatal logs a fatal error and exits the program
func fatal(ctx context.Context, msg string, err error) {
	if err != nil {
		slog.Error(msg, "error", err)
	} else {
		slog.Error(msg)
	}
	os.Exit(1)
}

// connectToMongoDB connects to MongoDB with retry logic
func connectToMongoDB(ctx context.Context, uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	slog.Info("Successfully connected to MongoDB")
	return client, nil
}
