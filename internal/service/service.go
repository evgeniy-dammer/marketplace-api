package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
)

// Authorization interface
type Authorization interface {
	GenerateToken(id string, username string, password string) (model.User, model.Tokens, error)
	ParseToken(token string) (string, error)
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

// Service interface
type Service struct {
	Authorization
	User
	Organization
	Category
	Item
}

// NewService constructor for Service
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User),
		Organization:  NewOrganizationService(repos.Organization),
		Category:      NewCategoryService(repos.Category),
		Item:          NewItemService(repos.Item),
	}
}
