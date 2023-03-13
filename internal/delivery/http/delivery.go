package http

import (
	"github.com/casbin/casbin-pg-adapter"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase"
)

// @title Emenu API
// @version 1.0
// @description Emenu API service
// @contact.name Evgeniy Dammer
// @contact.email evgeniydammer@gmail.com
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// Delivery delivery.
type Delivery struct {
	ucAuthentication usecase.Authentication
	ucAuthorization  usecase.Authorization
	ucUser           usecase.User
	ucOrganization   usecase.Organization
	ucCategory       usecase.Category
	ucItem           usecase.Item
	ucTable          usecase.Table
	ucOrder          usecase.Order
	ucImage          usecase.Image
	ucComment        usecase.Comment
	ucSpecification  usecase.Specification
	ucFavorite       usecase.Favorite
	ucRule           usecase.Rule
	adapter          *pgadapter.Adapter
	isTracingOn      bool
}

// New constructor for Delivery.
func New(
	ucAuthentication usecase.Authentication,
	ucAuthorization usecase.Authorization,
	ucUser usecase.User,
	ucOrganization usecase.Organization,
	ucCategory usecase.Category,
	ucItem usecase.Item,
	ucTable usecase.Table,
	ucOrder usecase.Order,
	ucImage usecase.Image,
	ucComment usecase.Comment,
	ucSpecification usecase.Specification,
	ucFavorite usecase.Favorite,
	ucRule usecase.Rule,
	adapter *pgadapter.Adapter,
	isTracingOn bool,
) *Delivery {
	return &Delivery{
		ucAuthentication: ucAuthentication,
		ucAuthorization:  ucAuthorization,
		ucUser:           ucUser,
		ucOrganization:   ucOrganization,
		ucCategory:       ucCategory,
		ucItem:           ucItem,
		ucTable:          ucTable,
		ucOrder:          ucOrder,
		ucImage:          ucImage,
		ucComment:        ucComment,
		ucSpecification:  ucSpecification,
		ucFavorite:       ucFavorite,
		ucRule:           ucRule,
		adapter:          adapter,
		isTracingOn:      isTracingOn,
	}
}
