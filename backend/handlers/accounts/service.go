// handlers/accounts/service.go
package accounts

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/Bedrockdude10/Booker/backend/domain"
	"github.com/Bedrockdude10/Booker/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	accounts *mongo.Collection
}

// NewService creates a new accounts service
func NewService(collections map[string]*mongo.Collection) *Service {
	return &Service{
		accounts: collections["accounts"],
	}
}

// Helper function to check if role is valid
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

//==============================================================================
// CreateAccount - Creates a new user account
//==============================================================================

func (s *Service) CreateAccount(ctx context.Context, params CreateAccountParams) (*Account, *utils.AppError) {
	// Validate role using the domain Set
	if !domain.ValidRoles.Has(params.Role) {
		return nil, utils.ValidationErrorLog(ctx, "Invalid role")
	}

	// Hash the password
	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		return nil, utils.InternalErrorLog(ctx, "Failed to hash password", err)
	}

	// Create account document
	account := Account{
		ID:           primitive.NewObjectID(),
		Email:        strings.ToLower(strings.TrimSpace(params.Email)), // Normalize email
		PasswordHash: hashedPassword,
		Role:         params.Role,
		Name:         strings.TrimSpace(params.Name),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true, // New accounts are active by default
	}

	// Insert into database
	if _, err := s.accounts.InsertOne(ctx, account); err != nil {
		// Check for duplicate key error (email already exists)
		if mongo.IsDuplicateKeyError(err) {
			return nil, utils.ValidationErrorLog(ctx, "An account with this email already exists")
		}
		return nil, utils.DatabaseErrorLog(ctx, "create account", err)
	}

	return &account, nil
}

//==============================================================================
// GetAccountByID - Retrieves account by ObjectID
//==============================================================================

func (s *Service) GetAccountByID(ctx context.Context, id primitive.ObjectID) (*Account, *utils.AppError) {
	// Validate ObjectID
	if id.IsZero() {
		return nil, utils.ValidationErrorLog(ctx, "Invalid account ID")
	}

	var account Account
	err := s.accounts.FindOne(ctx, bson.M{"_id": id}).Decode(&account)

	if err == mongo.ErrNoDocuments {
		return nil, utils.NotFoundLog(ctx, "Account")
	}
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find account by id", err)
	}

	return &account, nil
}

//==============================================================================
// GetAccountByEmail - Retrieves account by email address
//==============================================================================

func (s *Service) GetAccountByEmail(ctx context.Context, email string) (*Account, *utils.AppError) {
	// Validate email input
	if strings.TrimSpace(email) == "" {
		return nil, utils.ValidationErrorLog(ctx, "Email is required")
	}

	// Basic email format validation
	if !isValidEmail(email) {
		return nil, utils.ValidationErrorLog(ctx, "Invalid email format")
	}

	// Normalize email for search
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	var account Account
	err := s.accounts.FindOne(ctx, bson.M{"email": normalizedEmail}).Decode(&account)

	if err == mongo.ErrNoDocuments {
		return nil, utils.NotFoundLog(ctx, "Account")
	}
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "find account by email", err)
	}

	return &account, nil
}

//==============================================================================
// UpdateAccount - Updates account information
//==============================================================================

func (s *Service) UpdateAccount(ctx context.Context, id primitive.ObjectID, params UpdateAccountParams) (*Account, *utils.AppError) {
	// Validate ObjectID
	if id.IsZero() {
		return nil, utils.ValidationErrorLog(ctx, "Invalid account ID")
	}

	// Validate role if provided
	if params.Role != "" {
		// Validate role using the domain Set
		if !domain.ValidRoles.Has(params.Role) {
			return nil, utils.ValidationErrorLog(ctx, "Invalid role")
		}
	}

	// Build update document dynamically based on provided fields
	updateFields := bson.M{
		"updatedAt": time.Now(), // Always update the timestamp
	}

	// Only update fields that are provided
	if params.Email != "" {
		normalizedEmail := strings.ToLower(strings.TrimSpace(params.Email))
		if !isValidEmail(normalizedEmail) {
			return nil, utils.ValidationErrorLog(ctx, "Invalid email format")
		}
		updateFields["email"] = normalizedEmail
	}

	if params.Role != "" {
		updateFields["role"] = params.Role
	}

	if params.Name != "" {
		updateFields["name"] = strings.TrimSpace(params.Name)
	}

	// Use FindOneAndUpdate to get the updated document back
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedAccount Account
	err := s.accounts.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updateFields},
		opts,
	).Decode(&updatedAccount)

	if err == mongo.ErrNoDocuments {
		return nil, utils.NotFoundLog(ctx, "Account")
	}
	if err != nil {
		// Check for duplicate key error (email already exists)
		if mongo.IsDuplicateKeyError(err) {
			return nil, utils.ValidationErrorLog(ctx, "An account with this email already exists")
		}
		return nil, utils.DatabaseErrorLog(ctx, "update account", err)
	}

	return &updatedAccount, nil
}

//==============================================================================
// Authentication helper methods
//==============================================================================

// VerifyPassword verifies the password for the given email and returns the corresponding account if successful.
func (s *Service) VerifyPassword(ctx context.Context, email, password string) (*Account, *utils.AppError) {
	// Get account by email
	account, err := s.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err // Pass through the error (not found, validation, etc.)
	}

	// Check if account is active
	if !account.IsActive {
		return nil, utils.ValidationErrorLog(ctx, "Account is disabled")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)); err != nil {
		return nil, utils.ValidationErrorLog(ctx, "Invalid credentials")
	}

	return account, nil
}

// UpdatePassword updates a user's password (useful for password reset)
func (s *Service) UpdatePassword(ctx context.Context, id primitive.ObjectID, newPassword string) *utils.AppError {
	// Validate ObjectID
	if id.IsZero() {
		return utils.ValidationErrorLog(ctx, "Invalid account ID")
	}

	// Hash new password
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return utils.InternalErrorLog(ctx, "Failed to hash password", err)
	}

	// Update password in database
	result, err := s.accounts.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"passwordHash": hashedPassword,
			"updatedAt":    time.Now(),
		}},
	)

	if err != nil {
		return utils.DatabaseErrorLog(ctx, "update password", err)
	}

	if result.MatchedCount == 0 {
		return utils.NotFoundLog(ctx, "Account")
	}

	return nil
}

// DeactivateAccount sets IsActive to false (soft delete)
func (s *Service) DeactivateAccount(ctx context.Context, id primitive.ObjectID) *utils.AppError {
	// Validate ObjectID
	if id.IsZero() {
		return utils.ValidationErrorLog(ctx, "Invalid account ID")
	}

	result, err := s.accounts.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"isActive":  false,
			"updatedAt": time.Now(),
		}},
	)

	if err != nil {
		return utils.DatabaseErrorLog(ctx, "deactivate account", err)
	}

	if result.MatchedCount == 0 {
		return utils.NotFoundLog(ctx, "Account")
	}

	return nil
}

//==============================================================================
// Password hashing utilities
//==============================================================================

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	// Use bcrypt default cost (currently 10)
	// This provides a good balance of security and performance
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPasswordHash checks if a password matches a hash
func VerifyPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//==============================================================================
// Email validation utility
//==============================================================================

// Simple email validation using regex
func isValidEmail(email string) bool {
	// Basic email validation - you might want to use a more sophisticated library
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

//==============================================================================
// Additional helper methods for handlers
//==============================================================================

// GetActiveAccountByEmail - Only returns active accounts
func (s *Service) GetActiveAccountByEmail(ctx context.Context, email string) (*Account, *utils.AppError) {
	account, err := s.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !account.IsActive {
		return nil, utils.NotFoundLog(ctx, "Account") // Don't reveal that account exists but is inactive
	}

	return account, nil
}

// ListAccounts - For admin purposes (with pagination)
func (s *Service) ListAccounts(ctx context.Context, page, limit int) ([]Account, *utils.AppError) {
	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Set up find options
	findOptions := options.Find().
		SetSort(bson.M{"createdAt": -1}). // Most recent first
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := s.accounts.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "list accounts", err)
	}
	defer cursor.Close(ctx)

	var accounts []Account
	if err := cursor.All(ctx, &accounts); err != nil {
		return nil, utils.DatabaseErrorLog(ctx, "decode accounts", err)
	}

	return accounts, nil
}

// CountAccounts - Get total number of accounts (for pagination)
func (s *Service) CountAccounts(ctx context.Context) (int64, *utils.AppError) {
	count, err := s.accounts.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, utils.DatabaseErrorLog(ctx, "count accounts", err)
	}
	return count, nil
}
