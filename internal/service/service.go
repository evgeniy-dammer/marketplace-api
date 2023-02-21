package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
)

// Authorization interface.
type Authorization interface {
	GenerateToken(id string, username string, password string) (model.User, model.Tokens, error)
	ParseToken(token string) (string, error)
	CreateUser(user model.User) (string, error)
	GetUserRole(id string) (string, error)
}

// User interface.
type User interface {
	GetAll(search string, status string, roleID string) ([]model.User, error)
	GetOne(userID string) (model.User, error)
	Create(userID string, user model.User) (string, error)
	Update(userID string, input model.UpdateUserInput) error
	Delete(userID string, dUserID string) error

	GetAllRoles() ([]model.Role, error)
}

// Organization interface.
type Organization interface {
	GetAll(userID string) ([]model.Organization, error)
	GetOne(userID string, organizationID string) (model.Organization, error)
	Create(userID string, organization model.Organization) (string, error)
	Update(userID string, input model.UpdateOrganizationInput) error
	Delete(userID string, organizationID string) error
}

// Category interface.
type Category interface {
	GetAll(userID string, organizationID string) ([]model.Category, error)
	GetOne(userID string, organizationID string, categoryID string) (model.Category, error)
	Create(userID string, category model.Category) (string, error)
	Update(userID string, input model.UpdateCategoryInput) error
	Delete(userID string, organizationID string, categoryID string) error
}

// Item interface.
type Item interface {
	GetAll(userID string, organizationID string) ([]model.Item, error)
	GetOne(userID string, organizationID string, itemID string) (model.Item, error)
	Create(userID string, item model.Item) (string, error)
	Update(userID string, input model.UpdateItemInput) error
	Delete(userID string, organizationID string, itemID string) error
}

// Table interface.
type Table interface {
	GetAll(userID string, organizationID string) ([]model.Table, error)
	GetOne(userID string, organizationID string, tableID string) (model.Table, error)
	Create(userID string, table model.Table) (string, error)
	Update(userID string, input model.UpdateTableInput) error
	Delete(userID string, organizationID string, tableID string) error
}

// Order interface.
type Order interface {
	GetAll(userID string, organizationID string) ([]model.Order, error)
	GetOne(userID string, organizationID string, orderID string) (model.Order, error)
	Create(userID string, order model.Order) (string, error)
	Update(userID string, input model.UpdateOrderInput) error
	Delete(userID string, organizationID string, orderID string) error
}

// Image interface.
type Image interface {
	GetAll(userID string, organizationID string) ([]model.Image, error)
	GetOne(userID string, organizationID string, imageID string) (model.Image, error)
	Create(userID string, image model.Image) (string, error)
	Update(userID string, input model.UpdateImageInput) error
	Delete(userID string, organizationID string, imageID string) error
}

// Service interface.
type Service struct {
	Authorization
	User
	Organization
	Category
	Item
	Table
	Order
	Image
}

// NewService constructor for Service.
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User),
		Organization:  NewOrganizationService(repos.Organization),
		Category:      NewCategoryService(repos.Category),
		Item:          NewItemService(repos.Item),
		Table:         NewTableService(repos.Table),
		Order:         NewOrderService(repos.Order),
		Image:         NewImageService(repos.Image),
	}
}