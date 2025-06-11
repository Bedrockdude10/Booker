// handlers/accounts/auth.go
package accounts

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JWT Claims structure
type Claims struct {
	UserID primitive.ObjectID `json:"user_id"`
	Email  string             `json:"email"`
	Role   string             `json:"role"`
	Name   string             `json:"name"`
	jwt.RegisteredClaims
}

// JWT service for token operations
type JWTService struct {
	secretKey []byte
	issuer    string
}

// NewJWTService creates a new JWT service
func NewJWTService() *JWTService {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		// In production, this should be a required environment variable
		secretKey = "your-super-secret-jwt-key-change-this-in-production"
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "booker-api"
	}

	return &JWTService{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// GenerateToken creates a new JWT token for the given account
func (j *JWTService) GenerateToken(account *Account) (string, error) {
	// Token expires in 24 hours by default
	expirationTime := time.Now().Add(24 * time.Hour)

	// Allow custom expiration via environment variable
	if customDuration := os.Getenv("JWT_EXPIRATION_HOURS"); customDuration != "" {
		if hours, err := time.ParseDuration(customDuration + "h"); err == nil {
			expirationTime = time.Now().Add(hours)
		}
	}

	claims := &Claims{
		UserID: account.ID,
		Email:  account.Email,
		Role:   account.Role,
		Name:   account.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
			Subject:   account.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken validates and parses a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}

// RefreshToken creates a new token from an existing valid token
func (j *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Create new token with extended expiration
	newExpirationTime := time.Now().Add(24 * time.Hour)
	claims.ExpiresAt = jwt.NewNumericDate(newExpirationTime)
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}
