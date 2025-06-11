// handlers/accounts/types.go
package accounts

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"` // Changed: json:"id" instead of json:"_id,omitempty"
	Email        string             `bson:"email" json:"email" validate:"required,email"`
	PasswordHash string             `bson:"passwordHash" json:"-"` // Never return in JSON
	Role         string             `bson:"role" json:"role" validate:"required,validrole"`
	Name         string             `bson:"name" json:"name" validate:"required,min=1,max=100"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	IsActive     bool               `bson:"isActive" json:"isActive"`
}

type CreateAccountParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,validrole"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateAccountParams struct {
	Email string `json:"email,omitempty" validate:"omitempty,email"`
	Role  string `json:"role" validate:"required,validrole"`
	Name  string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
}
