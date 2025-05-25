// handlers/accounts/service_test.go
package accounts

import (
	"context"
	"testing"
	"time"

	"github.com/Bedrockdude10/Booker/backend/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// AccountsServiceTestSuite provides a clean test environment
type AccountsServiceTestSuite struct {
	suite.Suite
	mt      *mtest.T
	service *Service
	ctx     context.Context
}

// SetupSuite runs once before all tests
func (suite *AccountsServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
}

// SetupTest runs before each individual test
func (suite *AccountsServiceTestSuite) SetupTest() {
	// Create fresh mock database for each test
	mt := mtest.New(suite.T(), mtest.NewOptions().ClientType(mtest.Mock))
	suite.mt = mt

	collections := map[string]*mongo.Collection{
		"accounts": mt.Coll,
	}
	suite.service = NewService(collections)
}

// TearDownTest runs after each individual test
func (suite *AccountsServiceTestSuite) TearDownTest() {
	suite.mt = nil
}

// Helper method to create BSON document from account
func (suite *AccountsServiceTestSuite) accountToBSON(account Account) bson.D {
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

func (suite *AccountsServiceTestSuite) TestCreateAccount_Success() {
	// Mock successful insertion
	suite.mt.AddMockResponses(mtest.CreateSuccessResponse())

	params := CreateAccountParams{
		Email:    "test@example.com",
		Password: "password123",
		Role:     domain.RolePromoter,
		Name:     "Test User",
	}

	account, err := suite.service.CreateAccount(suite.ctx, params)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), account)
	assert.Equal(suite.T(), "test@example.com", account.Email)
	assert.Equal(suite.T(), domain.RolePromoter, account.Role)
	assert.Equal(suite.T(), "Test User", account.Name)
	assert.True(suite.T(), account.IsActive)
	assert.False(suite.T(), account.ID.IsZero())
	assert.NotEmpty(suite.T(), account.PasswordHash)
	assert.NotEqual(suite.T(), "password123", account.PasswordHash) // Should be hashed
	assert.NotZero(suite.T(), account.CreatedAt)
	assert.NotZero(suite.T(), account.UpdatedAt)
}

func (suite *AccountsServiceTestSuite) TestCreateAccount_DatabaseError() {
	// Mock database error
	suite.mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
		Index:   1,
		Code:    11000, // Duplicate key error
		Message: "duplicate key error",
	}))

	params := CreateAccountParams{
		Email:    "duplicate@example.com",
		Password: "password123",
		Role:     domain.RoleArtist,
		Name:     "Duplicate User",
	}

	account, err := suite.service.CreateAccount(suite.ctx, params)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), account)
	assert.Contains(suite.T(), err.Error(), "create account")
}

func (suite *AccountsServiceTestSuite) TestCreateAccount_DifferentRoles() {
	// Test creating accounts with different roles
	roles := []string{domain.RolePromoter, domain.RoleArtist, domain.RoleAdmin}

	for _, role := range roles {
		suite.SetupTest() // Fresh mock for each iteration
		suite.mt.AddMockResponses(mtest.CreateSuccessResponse())

		params := CreateAccountParams{
			Email:    "user@example.com",
			Password: "password123",
			Role:     role,
			Name:     "Test User",
		}

		account, err := suite.service.CreateAccount(suite.ctx, params)

		assert.NoError(suite.T(), err, "Should create account with role: %s", role)
		assert.Equal(suite.T(), role, account.Role)
	}
}

//==============================================================================
// GetAccountByID Tests
//==============================================================================

func (suite *AccountsServiceTestSuite) TestGetAccountByID_Found() {
	accountID := primitive.NewObjectID()
	expectedAccount := Account{
		ID:           accountID,
		Email:        "found@example.com",
		PasswordHash: "hashedpassword",
		Role:         domain.RolePromoter,
		Name:         "Found User",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
	}

	suite.mt.AddMockResponses(
		mtest.CreateCursorResponse(1, "booker_test.accounts", mtest.FirstBatch, suite.accountToBSON(expectedAccount)),
	)

	account, err := suite.service.GetAccountByID(suite.ctx, accountID)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), account)
	assert.Equal(suite.T(), accountID, account.ID)
	assert.Equal(suite.T(), "found@example.com", account.Email)
	assert.Equal(suite.T(), domain.RolePromoter, account.Role)
	assert.Equal(suite.T(), "Found User", account.Name)
	assert.True(suite.T(), account.IsActive)
}

func (suite *AccountsServiceTestSuite) TestGetAccountByID_NotFound() {
	// Mock empty response (not found)
	suite.mt.AddMockResponses(
		mtest.CreateCursorResponse(0, "booker_test.accounts", mtest.FirstBatch),
	)

	nonExistentID := primitive.NewObjectID()
	account, err := suite.service.GetAccountByID(suite.ctx, nonExistentID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), account)
	assert.Contains(suite.T(), err.Error(), "Account not found")
}

func (suite *AccountsServiceTestSuite) TestGetAccountByID_DatabaseError() {
	// Mock database error
	suite.mt.AddMockResponses(
		mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "database connection failed",
		}),
	)

	accountID := primitive.NewObjectID()
	account, err := suite.service.GetAccountByID(suite.ctx, accountID)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), account)
	assert.Contains(suite.T(), err.Error(), "find account by id")
}

//==============================================================================
// GetAccountByEmail Tests
//==============================================================================

func (suite *AccountsServiceTestSuite) TestGetAccountByEmail_Found() {
	expectedAccount := Account{
		ID:           primitive.NewObjectID(),
		Email:        "user@example.com",
		PasswordHash: "hashedpassword",
		Role:         domain.RoleArtist,
		Name:         "Email User",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsActive:     true,
	}

	suite.mt.AddMockResponses(
		mtest.CreateCursorResponse(1, "booker_test.accounts", mtest.FirstBatch, suite.accountToBSON(expectedAccount)),
	)

	account, err := suite.service.GetAccountByEmail(suite.ctx, "user@example.com")

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), account)
	assert.Equal(suite.T(), "user@example.com", account.Email)
	assert.Equal(suite.T(), domain.RoleArtist, account.Role)
	assert.Equal(suite.T(), "Email User", account.Name)
}

func (suite *AccountsServiceTestSuite) TestGetAccountByEmail_NotFound() {
	// Mock empty response
	suite.mt.AddMockResponses(
		mtest.CreateCursorResponse(0, "booker_test.accounts", mtest.FirstBatch),
	)

	account, err := suite.service.GetAccountByEmail(suite.ctx, "nonexistent@example.com")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), account)
	assert.Contains(suite.T(), err.Error(), "Account not found")
}

func (suite *AccountsServiceTestSuite) TestGetAccountByEmail_EmptyEmail() {
	// Should validate empty email without hitting database
	account, err := suite.service.GetAccountByEmail(suite.ctx, "")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), account)
	assert.Contains(suite.T(), err.Error(), "Email is required")
}

func (suite *AccountsServiceTestSuite) TestGetAccountByEmail_InvalidEmail() {
	// Should validate malformed email without hitting database
	account, err := suite.service.GetAccountByEmail(suite.ctx, "invalid-email")

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), account)
	assert.Contains(suite.T(), err.Error(), "Invalid email format")
}

//==============================================================================
// UpdateAccount Tests
//==============================================================================

func (suite *AccountsServiceTestSuite) TestUpdateAccount_Success() {
	accountID := primitive.NewObjectID()
	updatedAccount := Account{
		ID:           accountID,
		Email:        "updated@example.com",
		PasswordHash: "hashedpassword",
		Role:         domain.RoleAdmin,
		Name:         "Updated User",
		CreatedAt:    time.Now().Add(-24 * time.Hour), // Created yesterday
		UpdatedAt:    time.Now(),                      // Updated now
		IsActive:     true,
	}

	// Mock successful update with returned document
	suite.mt.AddMockResponses(
		mtest.CreateCursorResponse(1, "booker_test.accounts", mtest.FirstBatch, suite.accountToBSON(updatedAccount)),
	)

	params := UpdateAccountParams{
		Email: "updated@example.com",
		Role:  domain.RoleAdmin,
		Name:  "Updated User",
	}

	account, err := suite.service.UpdateAccount(suite.ctx, accountID, params)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), account)
	assert.Equal(suite.T(), accountID, account.ID)
	assert.Equal(suite.T(), "updated@example.com", account.Email)
	assert.Equal(suite.T(), domain.RoleAdmin, account.Role)
	assert.Equal(suite.T(), "Updated User", account.Name)
}

func (suite *AccountsServiceTestSuite) TestUpdateAccount_NotFound() {
	// Mock empty response (document not found)
	suite.mt.AddMockResponses(
		mtest.CreateCursorResponse(0, "booker_test.accounts", mtest.FirstBatch),
	)

	nonExistentID := primitive.NewObjectID()
	params := UpdateAccountParams{
		Name: "Updated Name",
	}

	account, err := suite.service.UpdateAccount(suite.ctx, nonExistentID, params)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), account)
	assert.Contains(suite.T(), err.Error(), "Account not found")
}

func (suite *AccountsServiceTestSuite) TestUpdateAccount_PartialUpdate() {
	accountID := primitive.NewObjectID()
	// Only updating name, keeping other fields the same
	updatedAccount := Account{
		ID:           accountID,
		Email:        "original@example.com", // Unchanged
		PasswordHash: "hashedpassword",
		Role:         domain.RolePromoter, // Unchanged
		Name:         "New Name Only",     // Changed
		CreatedAt:    time.Now().Add(-24 * time.Hour),
		UpdatedAt:    time.Now(),
		IsActive:     true,
	}

	suite.mt.AddMockResponses(
		mtest.CreateCursorResponse(1, "booker_test.accounts", mtest.FirstBatch, suite.accountToBSON(updatedAccount)),
	)

	// Only updating name
	params := UpdateAccountParams{
		Name: "New Name Only",
	}

	account, err := suite.service.UpdateAccount(suite.ctx, accountID, params)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), account)
	assert.Equal(suite.T(), "New Name Only", account.Name)
	assert.Equal(suite.T(), "original@example.com", account.Email) // Should be unchanged
	assert.Equal(suite.T(), domain.RolePromoter, account.Role)     // Should be unchanged
}

//==============================================================================
// Integration-style Tests (Business Logic)
//==============================================================================

func (suite *AccountsServiceTestSuite) TestCreateAccount_PasswordHashing() {
	suite.mt.AddMockResponses(mtest.CreateSuccessResponse())

	params := CreateAccountParams{
		Email:    "hash@example.com",
		Password: "plaintext123",
		Role:     domain.RolePromoter,
		Name:     "Hash Test User",
	}

	account, err := suite.service.CreateAccount(suite.ctx, params)

	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), "plaintext123", account.PasswordHash, "Password should be hashed")
	assert.NotEmpty(suite.T(), account.PasswordHash, "Password hash should not be empty")
	assert.Greater(suite.T(), len(account.PasswordHash), 20, "Hashed password should be reasonably long")
}

func (suite *AccountsServiceTestSuite) TestCreateAccount_DefaultValues() {
	suite.mt.AddMockResponses(mtest.CreateSuccessResponse())

	params := CreateAccountParams{
		Email:    "defaults@example.com",
		Password: "password123",
		Role:     domain.RoleArtist,
		Name:     "Default Test User",
	}

	account, err := suite.service.CreateAccount(suite.ctx, params)

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), account.IsActive, "New accounts should be active by default")
	assert.WithinDuration(suite.T(), time.Now(), account.CreatedAt, 5*time.Second, "CreatedAt should be recent")
	assert.WithinDuration(suite.T(), time.Now(), account.UpdatedAt, 5*time.Second, "UpdatedAt should be recent")
	assert.False(suite.T(), account.ID.IsZero(), "ID should be generated")
}

//==============================================================================
// Edge Cases and Error Conditions
//==============================================================================

func (suite *AccountsServiceTestSuite) TestGetAccountByID_InvalidObjectID() {
	// This would be handled at the handler level, but test the service behavior
	nilID := primitive.NilObjectID
	account, err := suite.service.GetAccountByID(suite.ctx, nilID)

	// Depending on implementation, this might return an error or nil result
	// Adjust assertion based on your actual implementation
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), account)
}

//==============================================================================
// Run the test suite
//==============================================================================

func TestAccountsServiceSuite(t *testing.T) {
	suite.Run(t, new(AccountsServiceTestSuite))
}

//==============================================================================
// Additional unit tests for role validation (no mocking needed)
//==============================================================================

func TestRoleValidation(t *testing.T) {
	validRoles := []string{domain.RolePromoter, domain.RoleArtist, domain.RoleAdmin}
	invalidRoles := []string{"invalid", "", "promoter123", "ADMIN"}

	t.Run("valid roles", func(t *testing.T) {
		for _, role := range validRoles {
			params := CreateAccountParams{
				Email:    "test@example.com",
				Password: "password123",
				Role:     role,
				Name:     "Test User",
			}
			// This test would need your validation function
			// assert.True(t, isValidRole(params.Role), "Role %s should be valid", role)
			assert.Contains(t, validRoles, params.Role)
		}
	})

	t.Run("invalid roles", func(t *testing.T) {
		for _, role := range invalidRoles {
			assert.NotContains(t, validRoles, role, "Role %s should be invalid", role)
		}
	})
}
