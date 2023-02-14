package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
)

// Authorization interface.
type Authorization interface {
	GenerateToken(id string, username string, password string) (model.User, model.Tokens, error)
	ParseToken(token string) (string, error)
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

// Service interface.
type Service struct {
	Authorization
	User
	Organization
	Category
	Item
}

// NewService constructor for Service.
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User),
		Organization:  NewOrganizationService(repos.Organization),
		Category:      NewCategoryService(repos.Category),
		Item:          NewItemService(repos.Item),
	}
}
