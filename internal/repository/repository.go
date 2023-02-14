package repository

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
)

// Authorization interface.
type Authorization interface {
	GetUser(id string, username string) (model.User, error)
	CreateUser(user model.User, statusID string) (string, error)
	GetUserRole(id string) (string, error)
}

// User interface.
type User interface {
	GetAll(search string, status string, roleID string) ([]model.User, error)
	GetOne(userID string) (model.User, error)
	Create(user model.User, statusID string) (string, error)
	Update(userID string, input model.UpdateUserInput) error
	Delete(userID string) error

	GetAllRoles() ([]model.Role, error)
	GetActiveStatusID(name string) (string, error)
}

// Organization interface.
type Organization interface {
	GetAll(userID string) ([]model.Organization, error)
	GetOne(userID string, organizationID string) (model.Organization, error)
	Create(userID string, organization model.Organization) (string, error)
	Update(userID string, organizationID string, input model.UpdateOrganizationInput) error
	Delete(userID string, organizationID string) error
}

// Category interface.
type Category interface {
	GetAll(userID string, organizationID string) ([]model.Category, error)
	GetOne(userID string, organizationID string, categoryID string) (model.Category, error)
	Create(userID string, organizationID string, category model.Category) (string, error)
	Update(userID string, organizationID string, categoryID string, input model.UpdateCategoryInput) error
	Delete(userID string, organizationID string, categoryID string) error
}

// Item interface.
type Item interface {
	GetAll(userID string, organizationID string) ([]model.Item, error)
	GetOne(userID string, organizationID string, itemID string) (model.Item, error)
	Create(userID string, organizationID string, item model.Item) (string, error)
	Update(userID string, organizationID string, itemID string, input model.UpdateItemInput) error
	Delete(userID string, organizationID string, itemID string) error
}

// Repository interface.
type Repository struct {
	Authorization
	User
	Organization
	Category
	Item
}

// NewRepository constructor for Repository.
func NewRepository(database *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(database),
		User:          NewUserPostgresql(database),
		Organization:  NewOrganizationPostgresql(database),
		Category:      NewCategoryPostgresql(database),
		Item:          NewItemPostgresql(database),
	}
}
