package cache

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/category"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/comment"
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

// Cache interface.
type Cache struct {
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
type Authentication interface{}

// Authorization interface.
type Authorization interface {
	AuthorizationGetUserRole(ctx context.Context, userID string) (string, error)
	AuthorizationSetUserRole(ctx context.Context, userID string, role string) error
}

// User interface.
type User interface {
	UserGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]user.User, error)
	UserSetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, users []user.User) error
	UserGetOne(ctx context.Context, userID string) (user.User, error)
	UserCreate(ctx context.Context, user user.User) error
	UserUpdate(ctx context.Context, user user.User) error
	UserDelete(ctx context.Context, userID string) error
	UserInvalidate(ctx context.Context) error

	UserGetAllRoles(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]role.Role, error)
	UserSetAllRoles(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, roles []role.Role) error
}

// Organization interface.
type Organization interface {
	OrganizationGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]organization.Organization, error)
	OrganizationSetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, organizations []organization.Organization) error
	OrganizationGetOne(ctx context.Context, organizationID string) (organization.Organization, error)
	OrganizationCreate(ctx context.Context, organization organization.Organization) error
	OrganizationUpdate(ctx context.Context, organization organization.Organization) error
	OrganizationDelete(ctx context.Context, organizationID string) error
	OrganizationInvalidate(ctx context.Context) error
}

// Category interface.
type Category interface {
	CategoryGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]category.Category, error)
	CategorySetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, categories []category.Category) error
	CategoryGetOne(ctx context.Context, categoryID string) (category.Category, error)
	CategoryCreate(ctx context.Context, category category.Category) error
	CategoryUpdate(ctx context.Context, category category.Category) error
	CategoryDelete(ctx context.Context, categoryID string) error
	CategoryInvalidate(ctx context.Context) error
}

// Item interface.
type Item interface {
	ItemGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]item.Item, error)
	ItemSetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, items []item.Item) error
	ItemGetOne(ctx context.Context, itemID string) (item.Item, error)
	ItemCreate(ctx context.Context, item item.Item) error
	ItemUpdate(ctx context.Context, item item.Item) error
	ItemDelete(ctx context.Context, itemID string) error
	ItemInvalidate(ctx context.Context) error
}

// Table interface.
type Table interface {
	TableGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]table.Table, error)
	TableSetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, tables []table.Table) error
	TableGetOne(ctx context.Context, tableID string) (table.Table, error)
	TableCreate(ctx context.Context, table table.Table) error
	TableUpdate(ctx context.Context, table table.Table) error
	TableDelete(ctx context.Context, tableID string) error
	TableInvalidate(ctx context.Context) error
}

// Order interface.
type Order interface {
	OrderGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]order.Order, error)
	OrderSetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, orders []order.Order) error
	OrderGetOne(ctx context.Context, orderID string) (order.Order, error)
	OrderCreate(ctx context.Context, order order.Order) error
	OrderUpdate(ctx context.Context, order order.Order) error
	OrderDelete(ctx context.Context, orderID string) error
	OrderInvalidate(ctx context.Context) error
}

// Image interface.
type Image interface {
	ImageGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]image.Image, error)
	ImageSetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, images []image.Image) error
	ImageGetOne(ctx context.Context, imageID string) (image.Image, error)
	ImageCreate(ctx context.Context, image image.Image) error
	ImageUpdate(ctx context.Context, image image.Image) error
	ImageDelete(ctx context.Context, imageID string) error
	ImageInvalidate(ctx context.Context) error
}

// Comment interface.
type Comment interface {
	CommentGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]comment.Comment, error)
	CommentSetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, comments []comment.Comment) error
	CommentGetOne(ctx context.Context, commentID string) (comment.Comment, error)
	CommentCreate(ctx context.Context, comment comment.Comment) error
	CommentUpdate(ctx context.Context, comment comment.Comment) error
	CommentDelete(ctx context.Context, commentID string) error
	CommentInvalidate(ctx context.Context) error
}

// Specification interface.
type Specification interface {
	SpecificationGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]specification.Specification, error)
	SpecificationSetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, specifications []specification.Specification) error
	SpecificationGetOne(ctx context.Context, specificationID string) (specification.Specification, error)
	SpecificationCreate(ctx context.Context, specification specification.Specification) error
	SpecificationUpdate(ctx context.Context, specification specification.Specification) error
	SpecificationDelete(ctx context.Context, specificationID string) error
	SpecificationInvalidate(ctx context.Context) error
}

// Favorite interface.
type Favorite interface{}

// Rule interface.
type Rule interface {
	RuleGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]rule.Rule, error)
	RuleSetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter, rules []rule.Rule) error
	RuleGetOne(ctx context.Context, ruleID string) (rule.Rule, error)
	RuleCreate(ctx context.Context, rule rule.Rule) error
	RuleUpdate(ctx context.Context, rule rule.Rule) error
	RuleDelete(ctx context.Context, ruleID string) error
	RuleInvalidate(ctx context.Context) error
}
