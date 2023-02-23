package repository

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

// Authorization interface.
type Authorization interface {
	GetUser(id string, username string) (model.User, error)
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

// Comment interface.
type Comment interface {
	GetAll(userID string, organizationID string) ([]model.Comment, error)
	GetOne(userID string, organizationID string, commentID string) (model.Comment, error)
	Create(userID string, comment model.Comment) (string, error)
	Update(userID string, input model.UpdateCommentInput) error
	Delete(userID string, organizationID string, commentID string) error
}

// Specification interface.
type Specification interface {
	GetAll(userID string, organizationID string) ([]model.Specification, error)
	GetOne(userID string, organizationID string, specificationID string) (model.Specification, error)
	Create(userID string, specification model.Specification) (string, error)
	Update(userID string, input model.UpdateSpecificationInput) error
	Delete(userID string, organizationID string, specificationID string) error
}

// Favorite interface.
type Favorite interface {
	Create(userID string, favorite model.Favorite) error
	Delete(userID string, itemID string) error
}

// Rule interface.
type Rule interface {
	GetAll(userID string) ([]model.Rule, error)
	GetOne(userID string, ruleID string) (model.Rule, error)
	Create(userID string, rule model.Rule) (string, error)
	Update(userID string, input model.UpdateRuleInput) error
	Delete(userID string, ruleID string) error
}

// Repository interface.
type Repository struct {
	Authorization
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

// NewRepository constructor for Repository.
func NewRepository(database *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgres(database),
		User:          postgres.NewUserPostgresql(database),
		Organization:  postgres.NewOrganizationPostgresql(database),
		Category:      postgres.NewCategoryPostgresql(database),
		Item:          postgres.NewItemPostgresql(database),
		Table:         postgres.NewTablePostgresql(database),
		Order:         postgres.NewOrderPostgresql(database),
		Image:         postgres.NewImagePostgresql(database),
		Comment:       postgres.NewCommentPostgresql(database),
		Specification: postgres.NewSpecificationPostgresql(database),
		Favorite:      postgres.NewFavoritePostgresql(database),
		Rule:          postgres.NewRulePostgresql(database),
	}
}
