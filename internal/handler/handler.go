package handler

import (
	"github.com/evgeniy-dammer/emenu-api/internal/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// Handler handler.
type Handler struct {
	services *service.Service
}

// NewHandler constructor for Handler.
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// InitRoutes crete routes.
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(h.corsMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	auth := router.Group("/auth")
	{
		auth.POST("/signin", h.signIn)
		auth.POST("/signup", h.signUp)
		auth.POST("/refresh", h.refresh)
	}

	api := router.Group("/api", h.userIdentity)
	{
		version1 := api.Group("/v1")
		{
			users := version1.Group("/users")
			{
				users.GET("/", h.getAllUsers)
				users.GET("/:id", h.getUser)
				users.POST("/", h.createUser)
				users.PATCH("/", h.updateUser)
				users.DELETE("/:id", h.deleteUser)
				users.GET("/roles", h.getAllRoles)
			}

			organizations := version1.Group("/organizations")
			{
				organizations.GET("/", h.getOrganizations)
				organizations.GET("/:id", h.getOrganization)
				organizations.POST("/", h.createOrganization)
				organizations.PATCH("/", h.updateOrganization)
				organizations.DELETE("/:id", h.deleteOrganization)
			}

			categories := version1.Group("/categories")
			{
				categories.GET("/:org_id", h.getCategories)
				categories.GET("/:org_id/:id", h.getCategory)
				categories.POST("/", h.createCategory)
				categories.PATCH("/", h.updateCategory)
				categories.DELETE("/:org_id/:id", h.deleteCategory)
			}

			items := version1.Group("/items")
			{
				items.GET("/:org_id", h.getItems)
				items.GET("/:org_id/:id", h.getItem)
				items.POST("/", h.createItem)
				items.PATCH("/", h.updateItem)
				items.DELETE("/:org_id/:id", h.deleteItem)
			}

			tables := version1.Group("/tables")
			{
				tables.GET("/:org_id", h.getTables)
				tables.GET("/:org_id/:id", h.getTable)
				tables.POST("/", h.createTable)
				tables.PATCH("/", h.updateTable)
				tables.DELETE("/:org_id/:id", h.deleteTable)
			}

			orders := version1.Group("/orders")
			{
				orders.GET("/:org_id", h.getOrders)
				orders.GET("/:org_id/:id", h.getOrder)
				orders.POST("/", h.createOrder)
				orders.PATCH("/", h.updateOrder)
				orders.DELETE("/:org_id/:id", h.deleteOrder)
			}
		}
	}

	return router
}
