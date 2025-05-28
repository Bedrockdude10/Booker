package accounts

import (
	"context"
	"testing"
	"time"

	"github.com/Bedrockdude10/Booker/backend/domain"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// Test helper to create service with mock collection
func setupService(mt *mtest.T) *Service {
	collections := map[string]*mongo.Collection{
		"accounts": mt.Coll,
	}
	return NewService(collections)
}

// Test helper to convert Account to BSON
func toBSON(account Account) bson.D {
	return bson.D{
		{Key: "_id", Value: account.ID},
		{Key: "email", Value: account.Email},
		{Key: "passwordHash", Value: account.PasswordHash},
		{Key: "role", Value: account.Role},
		{Key: "name", Value: account.Name},
		{Key: "createdAt", Value: account.CreatedAt},
		{Key: "updatedAt", Value: account.UpdatedAt},
		{Key: "isActive", Value: account.IsActive},
	}
}

//==============================================================================
// CreateAccount Tests
//==============================================================================

func TestCreateAccount_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("create account successfully", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		service := setupService(mt)

		params := CreateAccountParams{
			Email:    "test@example.com",
			Password: "password123",
			Role:     domain.RolePromoter,
			Name:     "Test User",
		}

		account, err := service.CreateAccount(context.Background(), params)

		assert.Nil(mt, err)
		assert.NotNil(mt, account)
		assert.Equal(mt, "test@example.com", account.Email)
		assert.Equal(mt, domain.RolePromoter, account.Role)
		assert.Equal(mt, "Test User", account.Name)
		assert.True(mt, account.IsActive)
		assert.NotEmpty(mt, account.PasswordHash)
		assert.NotEqual(mt, "password123", account.PasswordHash)
		assert.False(mt, account.ID.IsZero())
		assert.NotZero(mt, account.CreatedAt)
		assert.NotZero(mt, account.UpdatedAt)
	})
}

func TestCreateAccount_InvalidRole(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("invalid role", func(mt *mtest.T) {
		service := setupService(mt)

		params := CreateAccountParams{
			Email:    "test@example.com",
			Password: "password123",
			Role:     "invalid-role",
			Name:     "Test User",
		}

		account, err := service.CreateAccount(context.Background(), params)

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "Invalid role")
	})
}

func TestCreateAccount_DuplicateEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("duplicate email", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		service := setupService(mt)

		params := CreateAccountParams{
			Email:    "duplicate@example.com",
			Password: "password123",
			Role:     domain.RolePromoter,
			Name:     "Test User",
		}

		account, err := service.CreateAccount(context.Background(), params)

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "An account with this email already exists")
	})
}

//==============================================================================
// GetAccountByID Tests
//==============================================================================

func TestGetAccountByID_Found(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("account found", func(mt *mtest.T) {
		accountID := primitive.NewObjectID()
		expectedAccount := Account{
			ID:           accountID,
			Email:        "found@example.com",
			PasswordHash: "hashed",
			Role:         domain.RolePromoter,
			Name:         "Found User",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			IsActive:     true,
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.accounts", mtest.FirstBatch, toBSON(expectedAccount)),
		)
		service := setupService(mt)

		account, err := service.GetAccountByID(context.Background(), accountID)

		assert.Nil(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, accountID, account.ID)
		assert.Equal(t, "found@example.com", account.Email)
	})
}

func TestGetAccountByID_NotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("account not found", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.accounts", mtest.FirstBatch),
		)
		service := setupService(mt)

		account, err := service.GetAccountByID(context.Background(), primitive.NewObjectID())

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "Account not found")
	})
}

func TestGetAccountByID_InvalidID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("invalid id", func(mt *mtest.T) {
		service := setupService(mt)

		account, err := service.GetAccountByID(context.Background(), primitive.NilObjectID)

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "Invalid account ID")
	})
}

//==============================================================================
// GetAccountByEmail Tests
//==============================================================================

func TestGetAccountByEmail_Found(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("account found by email", func(mt *mtest.T) {
		expectedAccount := Account{
			ID:           primitive.NewObjectID(),
			Email:        "user@example.com",
			PasswordHash: "hashed",
			Role:         domain.RoleArtist,
			Name:         "Email User",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			IsActive:     true,
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.accounts", mtest.FirstBatch, toBSON(expectedAccount)),
		)
		service := setupService(mt)

		account, err := service.GetAccountByEmail(context.Background(), "user@example.com")

		assert.Nil(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, "user@example.com", account.Email)
	})
}

func TestGetAccountByEmail_NotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("account not found by email", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.accounts", mtest.FirstBatch),
		)
		service := setupService(mt)

		account, err := service.GetAccountByEmail(context.Background(), "notfound@example.com")

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "Account not found")
	})
}

func TestGetAccountByEmail_InvalidEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("invalid email format", func(mt *mtest.T) {
		service := setupService(mt)

		account, err := service.GetAccountByEmail(context.Background(), "invalid-email")

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "Invalid email format")
	})
}

func TestGetAccountByEmail_EmptyEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("empty email", func(mt *mtest.T) {
		service := setupService(mt)

		account, err := service.GetAccountByEmail(context.Background(), "")

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "Email is required")
	})
}

// commented out, need real data persistence not mocks

// //==============================================================================
// // UpdateAccount Tests
// //==============================================================================

// func TestUpdateAccount_Success(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

// 	mt.Run("update account successfully", func(mt *mtest.T) {
// 		service := setupService(mt)
// 		accountID := primitive.NewObjectID()

// 		// Create the expected updated account result
// 		updatedAccount := Account{
// 			ID:           accountID,
// 			Email:        "updated@example.com",
// 			PasswordHash: "existing-hash",
// 			Role:         domain.RoleArtist,
// 			Name:         "Updated User",
// 			CreatedAt:    time.Now().Add(-24 * time.Hour), // Created yesterday
// 			UpdatedAt:    time.Now(),                      // Updated now
// 			IsActive:     true,
// 		}

// 		// Mock the FindOneAndUpdate operation
// 		// This simulates finding and updating an existing account
// 		mt.AddMockResponses(
// 			mtest.CreateSuccessResponse(toBSON(updatedAccount)...),
// 		)

// 		// Test the update operation
// 		params := UpdateAccountParams{
// 			Email: "updated@example.com",
// 			Role:  domain.RoleArtist,
// 			Name:  "Updated User",
// 		}

// 		account, err := service.UpdateAccount(context.Background(), accountID, params)

// 		assert.Nil(mt, err)
// 		assert.NotNil(mt, account)
// 		assert.Equal(mt, "updated@example.com", account.Email)
// 		assert.Equal(mt, domain.RoleArtist, account.Role)
// 		assert.Equal(mt, "Updated User", account.Name)
// 		assert.Equal(mt, accountID, account.ID)
// 	})
// }

// func TestUpdateAccount_NotFound(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

// 	mt.Run("account not found for update", func(mt *mtest.T) {
// 		mt.AddMockResponses(
// 			mtest.CreateCursorResponse(0, "test.accounts", mtest.FirstBatch),
// 		)
// 		service := setupService(mt)

// 		params := UpdateAccountParams{Name: "Updated Name"}
// 		account, err := service.UpdateAccount(context.Background(), primitive.NewObjectID(), params)

// 		assert.Error(t, err)
// 		assert.Nil(t, account)
// 		assert.Contains(t, err.Error(), "Account not found")
// 	})
// }

//==============================================================================
// VerifyPassword Tests
//==============================================================================

func TestVerifyPassword_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("password verification success", func(mt *mtest.T) {
		// Hash for "password123"
		hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
		expectedAccount := Account{
			ID:           primitive.NewObjectID(),
			Email:        "user@example.com",
			PasswordHash: hashedPassword,
			Role:         domain.RolePromoter,
			Name:         "Test User",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			IsActive:     true,
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.accounts", mtest.FirstBatch, toBSON(expectedAccount)),
		)
		service := setupService(mt)

		account, err := service.VerifyPassword(context.Background(), "user@example.com", "password123")

		assert.Nil(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, "user@example.com", account.Email)
	})
}

func TestVerifyPassword_WrongPassword(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("wrong password", func(mt *mtest.T) {
		hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
		expectedAccount := Account{
			ID:           primitive.NewObjectID(),
			Email:        "user@example.com",
			PasswordHash: hashedPassword,
			Role:         domain.RolePromoter,
			Name:         "Test User",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			IsActive:     true,
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.accounts", mtest.FirstBatch, toBSON(expectedAccount)),
		)
		service := setupService(mt)

		account, err := service.VerifyPassword(context.Background(), "user@example.com", "wrongpassword")

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "Invalid credentials")
	})
}

func TestVerifyPassword_InactiveAccount(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("inactive account", func(mt *mtest.T) {
		hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
		expectedAccount := Account{
			ID:           primitive.NewObjectID(),
			Email:        "user@example.com",
			PasswordHash: hashedPassword,
			Role:         domain.RolePromoter,
			Name:         "Test User",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			IsActive:     false, // Inactive account
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.accounts", mtest.FirstBatch, toBSON(expectedAccount)),
		)
		service := setupService(mt)

		account, err := service.VerifyPassword(context.Background(), "user@example.com", "password123")

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "Account is disabled")
	})
}

//==============================================================================
// Role Validation Tests (no mocking needed)
//==============================================================================

func TestRoleValidation(t *testing.T) {
	validRoles := []string{domain.RolePromoter, domain.RoleArtist, domain.RoleAdmin}
	invalidRoles := []string{"invalid", "", "PROMOTER", "user"}

	t.Run("valid roles", func(t *testing.T) {
		for _, role := range validRoles {
			assert.True(t, domain.HasRole(role), "Role %s should be valid", role)
		}
	})

	t.Run("invalid roles", func(t *testing.T) {
		for _, role := range invalidRoles {
			assert.False(t, domain.HasRole(role), "Role %s should be invalid", role)
		}
	})
}
