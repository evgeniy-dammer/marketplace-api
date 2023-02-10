package handler

import (
	"github.com/evgeniy-dammer/emenu-api/internal/service"
	"github.com/gin-gonic/gin"
)

// Handler
type Handler struct {
	services *service.Service
}

// NewHandler constructor for Handler
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// InitRoutes crete routes
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(h.corsMiddleware())

	auth := router.Group("/auth")
	{
		auth.POST("/signin", h.signIn)
		auth.POST("/signup", h.signUp)
		auth.POST("/refresh", h.refresh)
	}

	api := router.Group("/api", h.userIdentity)
	{
		v1 := api.Group("/v1")
		{
			users := v1.Group("/users")
			{
				users.GET("/", h.getAllUsers)
				users.GET("/:id", h.getUser)
				users.POST("/", h.createUser)
				users.PATCH("/:id", h.updateUser)
				users.DELETE("/:id", h.deleteUser)
				users.GET("/roles", h.getAllRoles)
			}

			organizations := v1.Group("/organizations")
			{
				organizations.GET("/", h.getOrganizations)
				organizations.GET("/:id", h.getOrganization)
				organizations.POST("/", h.createOrganization)
				organizations.PATCH("/:id", h.updateOrganization)
				organizations.DELETE("/:id", h.deleteOrganization)
			}

			categories := v1.Group("/categories")
			{
				categories.GET("/:org_id", h.getCategories)
				categories.GET("/:org_id/:id", h.getCategory)
				categories.POST("/:org_id", h.createCategory)
				categories.PATCH("/:org_id/:id", h.updateCategory)
				categories.DELETE("/:org_id/:id", h.deleteCategory)
			}

			items := v1.Group("/items")
			{
				items.GET("/:org_id", h.getItems)
				items.GET("/:org_id/:id", h.getItem)
				items.POST("/:org_id", h.createItem)
				items.PATCH("/:org_id/:id", h.updateItem)
				items.DELETE("/:org_id/:id", h.deleteItem)
			}

		}
	}

	return router
}
