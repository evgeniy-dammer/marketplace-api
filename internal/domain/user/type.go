package user

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

//easyjson:json
type ListUser []User

// User entity.
//
//easyjson:json
type User struct {
	// User ID
	ID string `json:"id" db:"id"`
	// Phone number
	Phone string `json:"phone" db:"phone" binding:"required"`
	// Password
	Password string `json:"password,omitempty" db:"password" binding:"required"`
	// First name
	FirstName string `json:"firstname" db:"first_name" binding:"required"`
	// Last Name
	LastName string `json:"lastname" db:"last_name" binding:"required"`
	// Role name
	RoleName string `json:"role,omitempty" db:"role"`
	// Users status
	Status string `json:"status" db:"status"`
	// Role ID
	RoleID int `json:"roleid,omitempty"`
}

// CreateUserInput entity.
//
//easyjson:json
type CreateUserInput struct {
	// Phone number
	Phone string `json:"phone" db:"phone" binding:"required"`
	// Password
	Password string `json:"password,omitempty" db:"password" binding:"required"`
	// First name
	FirstName string `json:"firstname" db:"first_name" binding:"required"`
	// Last Name
	LastName string `json:"lastname" db:"last_name" binding:"required"`
	// Role ID
	RoleID int `json:"roleid,omitempty"`
}

// SignInInput is an input data for signing in.
type SignInInput struct {
	// Users phone number
	Phone string `json:"phone" binding:"required"`
	// Users password
	Password string `json:"password" binding:"required"`
}

// UpdateUserInput is an input data for updating user entity.
//
//easyjson:json
type UpdateUserInput struct {
	// User ID
	ID *string `json:"id"`
	// First name
	FirstName *string `json:"firstname"`
	// Last Name
	LastName *string `json:"lastname"`
	// Password
	Password *string `json:"password"`
}

// Validate checks if update input is nil.
func (i UpdateUserInput) Validate() error {
	if i.ID == nil && i.FirstName == nil && i.LastName == nil {
		return ErrStructHasNoValues
	}

	return nil
}
