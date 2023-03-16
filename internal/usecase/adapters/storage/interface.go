package storage

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/category"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/comment"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/favorite"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/image"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/item"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/order"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/organization"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/role"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/rule"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/specification"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/table"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/user"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
)

// Storage interface.
type Storage struct {
	Authentication
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

// Authentication interface.
type Authentication interface {
	AuthenticationGetUser(ctx context.Context, id string, username string) (user.User, error)
	AuthenticationCreateUser(ctx context.Context, input user.CreateUserInput) (string, error)
}

// Authorization interface.
type Authorization interface {
	AuthorizationGetUserRole(ctx context.Context, id string) (string, error)
}

// User interface.
type User interface {
	UserGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]user.User, error)
	UserGetOne(ctx context.Context, meta query.MetaData, userID string) (user.User, error)
	UserCreate(ctx context.Context, meta query.MetaData, input user.CreateUserInput) (string, error)
	UserUpdate(ctx context.Context, meta query.MetaData, input user.UpdateUserInput) error
	UserDelete(ctx context.Context, meta query.MetaData, userID string) error

	UserGetAllRoles(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]role.Role, error)
}

// Organization interface.
type Organization interface {
	OrganizationGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]organization.Organization, error)
	OrganizationGetOne(ctx context.Context, meta query.MetaData, organizationID string) (organization.Organization, error)
	OrganizationCreate(ctx context.Context, meta query.MetaData, input organization.CreateOrganizationInput) (string, error)
	OrganizationUpdate(ctx context.Context, meta query.MetaData, input organization.UpdateOrganizationInput) error
	OrganizationDelete(ctx context.Context, meta query.MetaData, organizationID string) error
}

// Category interface.
type Category interface {
	CategoryGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]category.Category, error)
	CategoryGetOne(ctx context.Context, meta query.MetaData, categoryID string) (category.Category, error)
	CategoryCreate(ctx context.Context, meta query.MetaData, input category.CreateCategoryInput) (string, error)
	CategoryUpdate(ctx context.Context, meta query.MetaData, input category.UpdateCategoryInput) error
	CategoryDelete(ctx context.Context, meta query.MetaData, categoryID string) error
}

// Item interface.
type Item interface {
	ItemGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]item.Item, error)
	ItemGetOne(ctx context.Context, meta query.MetaData, itemID string) (item.Item, error)
	ItemCreate(ctx context.Context, meta query.MetaData, input item.CreateItemInput) (string, error)
	ItemUpdate(ctx context.Context, meta query.MetaData, input item.UpdateItemInput) error
	ItemDelete(ctx context.Context, meta query.MetaData, itemID string) error
}

// Table interface.
type Table interface {
	TableGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]table.Table, error)
	TableGetOne(ctx context.Context, meta query.MetaData, tableID string) (table.Table, error)
	TableCreate(ctx context.Context, meta query.MetaData, input table.CreateTableInput) (string, error)
	TableUpdate(ctx context.Context, meta query.MetaData, input table.UpdateTableInput) error
	TableDelete(ctx context.Context, meta query.MetaData, tableID string) error
}

// Order interface.
type Order interface {
	OrderGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]order.Order, error)
	OrderGetOne(ctx context.Context, meta query.MetaData, orderID string) (order.Order, error)
	OrderCreate(ctx context.Context, meta query.MetaData, input order.CreateOrderInput) (string, error)
	OrderUpdate(ctx context.Context, meta query.MetaData, input order.UpdateOrderInput) error
	OrderDelete(ctx context.Context, meta query.MetaData, orderID string) error
}

// Image interface.
type Image interface {
	ImageGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]image.Image, error)
	ImageGetOne(ctx context.Context, meta query.MetaData, imageID string) (image.Image, error)
	ImageCreate(ctx context.Context, meta query.MetaData, input image.CreateImageInput) (string, error)
	ImageUpdate(ctx context.Context, meta query.MetaData, input image.UpdateImageInput) error
	ImageDelete(ctx context.Context, meta query.MetaData, imageID string) error
}

// Comment interface.
type Comment interface {
	CommentGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]comment.Comment, error)
	CommentGetOne(ctx context.Context, meta query.MetaData, commentID string) (comment.Comment, error)
	CommentCreate(ctx context.Context, meta query.MetaData, input comment.CreateCommentInput) (string, error)
	CommentUpdate(ctx context.Context, meta query.MetaData, input comment.UpdateCommentInput) error
	CommentDelete(ctx context.Context, meta query.MetaData, commentID string) error
}

// Specification interface.
type Specification interface {
	SpecificationGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]specification.Specification, error)
	SpecificationGetOne(ctx context.Context, meta query.MetaData, specificationID string) (specification.Specification, error)
	SpecificationCreate(ctx context.Context, meta query.MetaData, input specification.CreateSpecificationInput) (string, error)
	SpecificationUpdate(ctx context.Context, meta query.MetaData, input specification.UpdateSpecificationInput) error
	SpecificationDelete(ctx context.Context, meta query.MetaData, specificationID string) error
}

// Favorite interface.
type Favorite interface {
	FavoriteCreate(ctx context.Context, meta query.MetaData, favorite favorite.Favorite) error
	FavoriteDelete(ctx context.Context, meta query.MetaData, itemID string) error
}

// Rule interface.
type Rule interface {
	RuleGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]rule.Rule, error)
	RuleGetOne(ctx context.Context, meta query.MetaData, ruleID string) (rule.Rule, error)
	RuleCreate(ctx context.Context, meta query.MetaData, input rule.CreateRuleInput) (string, error)
	RuleUpdate(ctx context.Context, meta query.MetaData, input rule.UpdateRuleInput) error
	RuleDelete(ctx context.Context, meta query.MetaData, ruleID string) error
}
