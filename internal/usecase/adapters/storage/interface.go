package storage

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/category"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/comment"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/favorite"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/image"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/item"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/order"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/organization"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/role"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/rule"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/specification"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/table"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
)

// Storage interface.
type Storage struct {
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
	AuthenticationGetUser(id string, username string) (user.User, error)
	AuthenticationCreateUser(user user.User) (string, error)
	AuthenticationGetUserRole(id string) (string, error)
}

// User interface.
type User interface {
	UserGetAll(search string, status string, roleID string) ([]user.User, error)
	UserGetOne(userID string) (user.User, error)
	UserCreate(userID string, user user.User) (string, error)
	UserUpdate(userID string, input user.UpdateUserInput) error
	UserDelete(userID string, dUserID string) error

	UserGetAllRoles() ([]role.Role, error)
}

// Organization interface.
type Organization interface {
	OrganizationGetAll(userID string) ([]organization.Organization, error)
	OrganizationGetOne(userID string, organizationID string) (organization.Organization, error)
	OrganizationCreate(userID string, organization organization.Organization) (string, error)
	OrganizationUpdate(userID string, input organization.UpdateOrganizationInput) error
	OrganizationDelete(userID string, organizationID string) error
}

// Category interface.
type Category interface {
	CategoryGetAll(userID string, organizationID string) ([]category.Category, error)
	CategoryGetOne(userID string, organizationID string, categoryID string) (category.Category, error)
	CategoryCreate(userID string, category category.Category) (string, error)
	CategoryUpdate(userID string, input category.UpdateCategoryInput) error
	CategoryDelete(userID string, organizationID string, categoryID string) error
}

// Item interface.
type Item interface {
	ItemGetAll(userID string, organizationID string) ([]item.Item, error)
	ItemGetOne(userID string, organizationID string, itemID string) (item.Item, error)
	ItemCreate(userID string, item item.Item) (string, error)
	ItemUpdate(userID string, input item.UpdateItemInput) error
	ItemDelete(userID string, organizationID string, itemID string) error
}

// Table interface.
type Table interface {
	TableGetAll(userID string, organizationID string) ([]table.Table, error)
	TableGetOne(userID string, organizationID string, tableID string) (table.Table, error)
	TableCreate(userID string, table table.Table) (string, error)
	TableUpdate(userID string, input table.UpdateTableInput) error
	TableDelete(userID string, organizationID string, tableID string) error
}

// Order interface.
type Order interface {
	OrderGetAll(userID string, organizationID string) ([]order.Order, error)
	OrderGetOne(userID string, organizationID string, orderID string) (order.Order, error)
	OrderCreate(userID string, order order.Order) (string, error)
	OrderUpdate(userID string, input order.UpdateOrderInput) error
	OrderDelete(userID string, organizationID string, orderID string) error
}

// Image interface.
type Image interface {
	ImageGetAll(userID string, organizationID string) ([]image.Image, error)
	ImageGetOne(userID string, organizationID string, imageID string) (image.Image, error)
	ImageCreate(userID string, image image.Image) (string, error)
	ImageUpdate(userID string, input image.UpdateImageInput) error
	ImageDelete(userID string, organizationID string, imageID string) error
}

// Comment interface.
type Comment interface {
	CommentGetAll(userID string, organizationID string) ([]comment.Comment, error)
	CommentGetOne(userID string, organizationID string, commentID string) (comment.Comment, error)
	CommentCreate(userID string, comment comment.Comment) (string, error)
	CommentUpdate(userID string, input comment.UpdateCommentInput) error
	CommentDelete(userID string, organizationID string, commentID string) error
}

// Specification interface.
type Specification interface {
	SpecificationGetAll(userID string, organizationID string) ([]specification.Specification, error)
	SpecificationGetOne(userID string, organizationID string, specificationID string) (specification.Specification, error)
	SpecificationCreate(userID string, specification specification.Specification) (string, error)
	SpecificationUpdate(userID string, input specification.UpdateSpecificationInput) error
	SpecificationDelete(userID string, organizationID string, specificationID string) error
}

// Favorite interface.
type Favorite interface {
	FavoriteCreate(userID string, favorite favorite.Favorite) error
	FavoriteDelete(userID string, itemID string) error
}

// Rule interface.
type Rule interface {
	RuleGetAll(userID string) ([]rule.Rule, error)
	RuleGetOne(userID string, ruleID string) (rule.Rule, error)
	RuleCreate(userID string, rule rule.Rule) (string, error)
	RuleUpdate(userID string, input rule.UpdateRuleInput) error
	RuleDelete(userID string, ruleID string) error
}
