package cache

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
)

// Cache interface.
type Cache struct {
	Authentication
	User
	Organization
	Category
	Item
	Table
	Order
	Image
	Comment
	Specification
	Favorite
	Rule
}

// Authentication interface.
type Authentication interface {
	AuthenticationGetUserRole(ctx context.Context, userID string) (string, error)
	AuthenticationSetUserRole(ctx context.Context, userID string, role string) error
}

// User interface.
type User interface {
	UserGetOne(ctx context.Context, userID string) (user.User, error)
	UserCreate(ctx context.Context, user user.User) error
	UserUpdate(ctx context.Context, user user.User) error
	UserDelete(ctx context.Context, userID string) error
}

// Organization interface.
type Organization interface{}

// Category interface.
type Category interface{}

// Item interface.
type Item interface{}

// Table interface.
type Table interface{}

// Order interface.
type Order interface{}

// Image interface.
type Image interface{}

// Comment interface.
type Comment interface{}

// Specification interface.
type Specification interface{}

// Favorite interface.
type Favorite interface{}

// Rule interface.
type Rule interface{}
