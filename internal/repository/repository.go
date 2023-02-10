package repository

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
)

// Authorization interface
type Authorization interface {
	GetUser(id string, username string) (model.User, error)
	CreateUser(user model.User, statusId string) (string, error)
	GetUserRole(id string) (string, error)
}

// User interface
type User interface {
	GetAll(search string, status string, roleId string) ([]model.User, error)
	GetOne(userId string) (model.User, error)
	Create(user model.User, statusId string) (string, error)
	Update(userId string, input model.UpdateUserInput) error
	Delete(userId string) error

	GetAllRoles() ([]model.Role, error)
	GetActiveStatusId(name string) (string, error)
}

// Organization interface
type Organization interface {
	GetAll(userId string) ([]model.Organization, error)
	GetOne(userId string, organizationId string) (model.Organization, error)
	Create(userId string, organization model.Organization) (string, error)
	Update(userId string, organizationId string, input model.UpdateOrganizationInput) error
	Delete(userId string, organizationId string) error
}

// Category interface
type Category interface {
	GetAll(userId string, organizationId string) ([]model.Category, error)
	GetOne(userId string, organizationId string, categoryId string) (model.Category, error)
	Create(userId string, organizationId string, category model.Category) (string, error)
	Update(userId string, organizationId string, categoryId string, input model.UpdateCategoryInput) error
	Delete(userId string, organizationId string, categoryId string) error
}

// Item interface
type Item interface {
	GetAll(userId string, organizationId string) ([]model.Item, error)
	GetOne(userId string, organizationId string, itemId string) (model.Item, error)
	Create(userId string, organizationId string, item model.Item) (string, error)
	Update(userId string, organizationId string, itemId string, input model.UpdateItemInput) error
	Delete(userId string, organizationId string, itemId string) error
}

// Repository interface
type Repository struct {
	Authorization
	User
	Organization
	Category
	Item
}

// NewRepository constructor for Repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		User:          NewUserPostgresql(db),
		Organization:  NewOrganizationPostgresql(db),
		Category:      NewCategoryPostgresql(db),
		Item:          NewItemPostgresql(db),
	}
}
