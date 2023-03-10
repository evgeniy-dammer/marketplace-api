package user

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// User entity.
type User struct {
	ID        string `json:"id" db:"id"`
	Phone     string `json:"phone" db:"phone" binding:"required"`
	Password  string `json:"password,omitempty" db:"password" binding:"required"`
	FirstName string `json:"firstname" db:"first_name" binding:"required"`
	LastName  string `json:"lastname" db:"last_name" binding:"required"`
	RoleName  string `json:"role,omitempty" db:"role"`
	Status    string `json:"status" db:"status"`
	RoleID    int    `json:"roleid,omitempty"`
}

// CreateUserInput entity.
type CreateUserInput struct {
	Phone     string `json:"phone" db:"phone" binding:"required"`
	Password  string `json:"password,omitempty" db:"password" binding:"required"`
	FirstName string `json:"firstname" db:"first_name" binding:"required"`
	LastName  string `json:"lastname" db:"last_name" binding:"required"`
	RoleID    int    `json:"roleid,omitempty"`
}

// SignInInput is an input data for signing in.
type SignInInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserInput is an input data for updating user entity.
type UpdateUserInput struct {
	ID        *string `json:"id"`
	FirstName *string `json:"firstname"`
	LastName  *string `json:"lastname"`
	Password  *string `json:"password"`
}

// Validate checks if update input is nil.
func (i UpdateUserInput) Validate() error {
	if i.ID == nil && i.FirstName == nil && i.LastName == nil {
		return ErrStructHasNoValues
	}

	return nil
}
