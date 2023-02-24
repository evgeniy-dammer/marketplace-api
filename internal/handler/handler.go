package handler

import (
	"github.com/casbin/casbin-pg-adapter"
	"github.com/evgeniy-dammer/emenu-api/internal/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// Handler handler.
type Handler struct {
	services *service.Service
	adapter  *pgadapter.Adapter
}

// NewHandler constructor for Handler.
func NewHandler(services *service.Service, adapter *pgadapter.Adapter) *Handler {
	return &Handler{services: services, adapter: adapter}
}

// InitRoutes crete routes.
func (h *Handler) InitRoutes(mode string) *gin.Engine {
	gin.SetMode(mode)

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
				users.GET("/", h.Authorize("users", "get", h.adapter), h.getAllUsers)
				users.GET("/:id", h.Authorize("user", "get", h.adapter), h.getUser)
				users.POST("/", h.Authorize("user", "post", h.adapter), h.createUser)
				users.PATCH("/", h.Authorize("user", "patch", h.adapter), h.updateUser)
				users.DELETE("/:id", h.Authorize("user", "delete", h.adapter), h.deleteUser)
				users.GET("/roles", h.Authorize("roles", "get", h.adapter), h.getAllRoles)
			}

			organizations := version1.Group("/organizations")
			{
				organizations.GET("/", h.Authorize("organizations", "get", h.adapter), h.getOrganizations)
				organizations.GET("/:id", h.Authorize("organization", "get", h.adapter), h.getOrganization)
				organizations.POST("/", h.Authorize("organization", "post", h.adapter), h.createOrganization)
				organizations.PATCH("/", h.Authorize("organization", "patch", h.adapter), h.updateOrganization)
				organizations.DELETE("/:id", h.Authorize("organization", "delete", h.adapter), h.deleteOrganization)
			}

			categories := version1.Group("/categories")
			{
				categories.GET("/:org_id", h.Authorize("categories", "get", h.adapter), h.getCategories)
				categories.GET("/:org_id/:id", h.Authorize("category", "get", h.adapter), h.getCategory)
				categories.POST("/", h.Authorize("category", "post", h.adapter), h.createCategory)
				categories.PATCH("/", h.Authorize("category", "patch", h.adapter), h.updateCategory)
				categories.DELETE("/:org_id/:id", h.Authorize("category", "delete", h.adapter), h.deleteCategory)
			}

			items := version1.Group("/items")
			{
				items.GET("/:org_id", h.Authorize("items", "get", h.adapter), h.getItems)
				items.GET("/:org_id/:id", h.Authorize("item", "get", h.adapter), h.getItem)
				items.POST("/", h.Authorize("item", "post", h.adapter), h.createItem)
				items.PATCH("/", h.Authorize("item", "patch", h.adapter), h.updateItem)
				items.DELETE("/:org_id/:id", h.Authorize("item", "delete", h.adapter), h.deleteItem)
			}

			tables := version1.Group("/tables")
			{
				tables.GET("/:org_id", h.Authorize("tables", "get", h.adapter), h.getTables)
				tables.GET("/:org_id/:id", h.Authorize("table", "get", h.adapter), h.getTable)
				tables.POST("/", h.Authorize("table", "post", h.adapter), h.createTable)
				tables.PATCH("/", h.Authorize("table", "patch", h.adapter), h.updateTable)
				tables.DELETE("/:org_id/:id", h.Authorize("table", "delete", h.adapter), h.deleteTable)
			}

			orders := version1.Group("/orders")
			{
				orders.GET("/:org_id", h.Authorize("orders", "get", h.adapter), h.getOrders)
				orders.GET("/:org_id/:id", h.Authorize("order", "get", h.adapter), h.getOrder)
				orders.POST("/", h.Authorize("order", "post", h.adapter), h.createOrder)
				orders.PATCH("/", h.Authorize("order", "patch", h.adapter), h.updateOrder)
				orders.DELETE("/:org_id/:id", h.Authorize("order", "delete", h.adapter), h.deleteOrder)
			}

			images := version1.Group("/images")
			{
				images.GET("/:org_id", h.Authorize("images", "get", h.adapter), h.getImages)
				images.GET("/:org_id/:id", h.Authorize("image", "get", h.adapter), h.getImage)
				images.POST("/", h.Authorize("image", "post", h.adapter), h.createImage)
				images.PATCH("/", h.Authorize("image", "patch", h.adapter), h.updateImage)
				images.DELETE("/:org_id/:id", h.Authorize("image", "delete", h.adapter), h.deleteImage)
			}

			comments := version1.Group("/comments")
			{
				comments.GET("/:org_id", h.Authorize("comments", "get", h.adapter), h.getComments)
				comments.GET("/:org_id/:id", h.Authorize("comment", "get", h.adapter), h.getComment)
				comments.POST("/", h.Authorize("comment", "post", h.adapter), h.createComment)
				comments.PATCH("/", h.Authorize("comment", "patch", h.adapter), h.updateComment)
				comments.DELETE("/:org_id/:id", h.Authorize("comment", "delete", h.adapter), h.deleteComment)
			}

			specifications := version1.Group("/specifications")
			{
				specifications.GET("/:org_id", h.Authorize("specifications", "get", h.adapter), h.getSpecifications)
				specifications.GET("/:org_id/:id", h.Authorize("specification", "get", h.adapter), h.getSpecification)
				specifications.POST("/", h.Authorize("specification", "post", h.adapter), h.createSpecification)
				specifications.PATCH("/", h.Authorize("specification", "patch", h.adapter), h.updateSpecification)
				specifications.DELETE("/:org_id/:id", h.Authorize("specification", "delete", h.adapter), h.deleteSpecification)
			}

			favorites := version1.Group("/favorites")
			{
				favorites.POST("/", h.Authorize("favorite", "post", h.adapter), h.createFavorite)
				favorites.DELETE("/:item_id", h.Authorize("favorite", "delete", h.adapter), h.deleteFavorite)
			}

			rules := version1.Group("/rules")
			{
				rules.GET("/", h.Authorize("rules", "get", h.adapter), h.getRules)
				rules.GET("/:id", h.Authorize("rule", "get", h.adapter), h.getRule)
				rules.POST("/", h.Authorize("rule", "post", h.adapter), h.createRule)
				rules.PATCH("/", h.Authorize("rule", "patch", h.adapter), h.updateRule)
				rules.DELETE("/:id", h.Authorize("rule", "delete", h.adapter), h.deleteRule)
			}
		}
	}

	return router
}
