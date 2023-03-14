package usecase

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
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/token"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/user"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
)

// Authentication interface.
type Authentication interface {
	AuthenticationGenerateToken(ctx context.Context, id string, username string, password string) (user.User, token.Tokens, error)
	AuthenticationParseToken(ctx context.Context, token string) (string, error)
	AuthenticationCreateUser(ctx context.Context, input user.CreateUserInput) (string, error)
}

// Authorization interface.
type Authorization interface {
	AuthorizationGetUserRole(ctx context.Context, id string) (string, error)
}

// User interface.
type User interface {
	UserGetAll(ctx context.Context, search string, status string, roleID string) ([]user.User, error)
	UserGetOne(ctx context.Context, userID string) (user.User, error)
	UserCreate(ctx context.Context, userID string, input user.CreateUserInput) (string, error)
	UserUpdate(ctx context.Context, userID string, input user.UpdateUserInput) error
	UserDelete(ctx context.Context, userID string, dUserID string) error

	UserGetAllRoles(ctx context.Context) ([]role.Role, error)
}

// Organization interface.
type Organization interface {
	OrganizationGetAll(ctx context.Context, userID string) ([]organization.Organization, error)
	OrganizationGetOne(ctx context.Context, userID string, organizationID string) (organization.Organization, error)
	OrganizationCreate(ctx context.Context, userID string, input organization.CreateOrganizationInput) (string, error)
	OrganizationUpdate(ctx context.Context, userID string, input organization.UpdateOrganizationInput) error
	OrganizationDelete(ctx context.Context, userID string, organizationID string) error
}

// Category interface.
type Category interface {
	CategoryGetAll(ctx context.Context, userID string, organizationID string) ([]category.Category, error)
	CategoryGetOne(ctx context.Context, userID string, organizationID string, categoryID string) (category.Category, error)
	CategoryCreate(ctx context.Context, userID string, input category.CreateCategoryInput) (string, error)
	CategoryUpdate(ctx context.Context, userID string, input category.UpdateCategoryInput) error
	CategoryDelete(ctx context.Context, userID string, organizationID string, categoryID string) error
}

// Item interface.
type Item interface {
	ItemGetAll(ctx context.Context, userID string, organizationID string) ([]item.Item, error)
	ItemGetOne(ctx context.Context, userID string, organizationID string, itemID string) (item.Item, error)
	ItemCreate(ctx context.Context, userID string, input item.CreateItemInput) (string, error)
	ItemUpdate(ctx context.Context, userID string, input item.UpdateItemInput) error
	ItemDelete(ctx context.Context, userID string, organizationID string, itemID string) error
}

// Table interface.
type Table interface {
	TableGetAll(ctx context.Context, userID string, organizationID string) ([]table.Table, error)
	TableGetOne(ctx context.Context, userID string, organizationID string, tableID string) (table.Table, error)
	TableCreate(ctx context.Context, userID string, input table.CreateTableInput) (string, error)
	TableUpdate(ctx context.Context, userID string, input table.UpdateTableInput) error
	TableDelete(ctx context.Context, userID string, organizationID string, tableID string) error
}

// Order interface.
type Order interface {
	OrderGetAll(ctx context.Context, userID string, organizationID string) ([]order.Order, error)
	OrderGetOne(ctx context.Context, userID string, organizationID string, orderID string) (order.Order, error)
	OrderCreate(ctx context.Context, userID string, input order.CreateOrderInput) (string, error)
	OrderUpdate(ctx context.Context, userID string, input order.UpdateOrderInput) error
	OrderDelete(ctx context.Context, userID string, organizationID string, orderID string) error
}

// Image interface.
type Image interface {
	ImageGetAll(ctx context.Context, userID string, organizationID string) ([]image.Image, error)
	ImageGetOne(ctx context.Context, userID string, organizationID string, imageID string) (image.Image, error)
	ImageCreate(ctx context.Context, userID string, input image.CreateImageInput) (string, error)
	ImageUpdate(ctx context.Context, userID string, input image.UpdateImageInput) error
	ImageDelete(ctx context.Context, userID string, organizationID string, imageID string) error
}

// Comment interface.
type Comment interface {
	CommentGetAll(ctx context.Context, userID string, organizationID string) ([]comment.Comment, error)
	CommentGetOne(ctx context.Context, userID string, organizationID string, commentID string) (comment.Comment, error)
	CommentCreate(ctx context.Context, userID string, input comment.CreateCommentInput) (string, error)
	CommentUpdate(ctx context.Context, userID string, input comment.UpdateCommentInput) error
	CommentDelete(ctx context.Context, userID string, organizationID string, commentID string) error
}

// Specification interface.
type Specification interface {
	SpecificationGetAll(ctx context.Context, userID string, organizationID string) ([]specification.Specification, error)
	SpecificationGetOne(ctx context.Context, userID string, organizationID string, specificationID string) (specification.Specification, error)
	SpecificationCreate(ctx context.Context, userID string, input specification.CreateSpecificationInput) (string, error)
	SpecificationUpdate(ctx context.Context, userID string, input specification.UpdateSpecificationInput) error
	SpecificationDelete(ctx context.Context, userID string, organizationID string, specificationID string) error
}

// Favorite interface.
type Favorite interface {
	FavoriteCreate(ctx context.Context, userID string, favorite favorite.Favorite) error
	FavoriteDelete(ctx context.Context, userID string, itemID string) error
}

// Rule interface.
type Rule interface {
	RuleGetAll(ctx context.Context, userID string) ([]rule.Rule, error)
	RuleGetOne(ctx context.Context, userID string, ruleID string) (rule.Rule, error)
	RuleCreate(ctx context.Context, userID string, input rule.CreateRuleInput) (string, error)
	RuleUpdate(ctx context.Context, userID string, input rule.UpdateRuleInput) error
	RuleDelete(ctx context.Context, userID string, ruleID string) error
}
