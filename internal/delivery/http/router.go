package http

import (
	"github.com/evgeniy-dammer/marketplace-api/docs"
	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// InitRoutes initialize routes.
func (d *Delivery) InitRoutes(mode string) *gin.Engine {
	gin.SetMode(mode)

	router := gin.New()
	router.Use(otelgin.Middleware("marketplace-api"))
	router.Use(d.corsMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(ginzap.RecoveryWithZap(logger.Logger, true))
	router.RedirectTrailingSlash = false

	if mode == "debug" {
		pprof.Register(router, "dev/pprof")
	}

	auth := router.Group("/auth")
	{
		auth.POST("/signin", d.signIn)
		auth.POST("/signup", d.signUp)
		auth.POST("/refresh", d.refresh)
	}

	metrics := router.Group("/metrics")
	{
		metrics.GET("/", gin.WrapH(promhttp.Handler()))
	}

	docs.SwaggerInfo_swagger.BasePath = "/" //nolint:nosnakecase

	router.Any("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api", d.userIdentity)
	{
		version1 := api.Group("/v1")
		{
			users := version1.Group("/users")
			{
				users.GET("", d.Authorize("users", "get", d.adapterSlave), d.getAllUsers)
				users.GET("/:id", d.Authorize("user", "get", d.adapterSlave), d.getUser)
				users.POST("", d.Authorize("user", "post", d.adapterSlave), d.createUser)
				users.PATCH("", d.Authorize("user", "patch", d.adapterSlave), d.updateUser)
				users.DELETE("/:id", d.Authorize("user", "delete", d.adapterSlave), d.deleteUser)
				users.GET("/roles", d.Authorize("roles", "get", d.adapterSlave), d.getAllRoles)
			}

			organizations := version1.Group("/organizations")
			{
				organizations.GET("", d.Authorize("organizations", "get", d.adapterSlave), d.getOrganizations)
				organizations.GET("/:id", d.Authorize("organization", "get", d.adapterSlave), d.getOrganization)
				organizations.POST("", d.Authorize("organization", "post", d.adapterSlave), d.createOrganization)
				organizations.PATCH("", d.Authorize("organization", "patch", d.adapterSlave), d.updateOrganization)
				organizations.DELETE("/:id", d.Authorize("organization", "delete", d.adapterSlave), d.deleteOrganization)
			}

			categories := version1.Group("/categories")
			{
				categories.GET("", d.Authorize("categories", "get", d.adapterSlave), d.getCategories)
				categories.GET("/:id", d.Authorize("category", "get", d.adapterSlave), d.getCategory)
				categories.POST("", d.Authorize("category", "post", d.adapterSlave), d.createCategory)
				categories.PATCH("", d.Authorize("category", "patch", d.adapterSlave), d.updateCategory)
				categories.DELETE("/:id", d.Authorize("category", "delete", d.adapterSlave), d.deleteCategory)
			}

			items := version1.Group("/items")
			{
				items.GET("", d.Authorize("items", "get", d.adapterSlave), d.getItems)
				items.GET("/:id", d.Authorize("item", "get", d.adapterSlave), d.getItem)
				items.POST("", d.Authorize("item", "post", d.adapterSlave), d.createItem)
				items.PATCH("", d.Authorize("item", "patch", d.adapterSlave), d.updateItem)
				items.DELETE("/:id", d.Authorize("item", "delete", d.adapterSlave), d.deleteItem)
			}

			tables := version1.Group("/tables")
			{
				tables.GET("", d.Authorize("tables", "get", d.adapterSlave), d.getTables)
				tables.GET("/:id", d.Authorize("table", "get", d.adapterSlave), d.getTable)
				tables.POST("", d.Authorize("table", "post", d.adapterSlave), d.createTable)
				tables.PATCH("", d.Authorize("table", "patch", d.adapterSlave), d.updateTable)
				tables.DELETE("/:id", d.Authorize("table", "delete", d.adapterSlave), d.deleteTable)
			}

			orders := version1.Group("/orders")
			{
				orders.GET("", d.Authorize("orders", "get", d.adapterSlave), d.getOrders)
				orders.GET("/:id", d.Authorize("order", "get", d.adapterSlave), d.getOrder)
				orders.POST("", d.Authorize("order", "post", d.adapterSlave), d.createOrder)
				orders.PATCH("", d.Authorize("order", "patch", d.adapterSlave), d.updateOrder)
				orders.DELETE("/:id", d.Authorize("order", "delete", d.adapterSlave), d.deleteOrder)
			}

			images := version1.Group("/images")
			{
				images.GET("", d.Authorize("images", "get", d.adapterSlave), d.getImages)
				images.GET("/:id", d.Authorize("image", "get", d.adapterSlave), d.getImage)
				images.POST("", d.Authorize("image", "post", d.adapterSlave), d.createImage)
				images.PATCH("", d.Authorize("image", "patch", d.adapterSlave), d.updateImage)
				images.DELETE("/:id", d.Authorize("image", "delete", d.adapterSlave), d.deleteImage)
			}

			comments := version1.Group("/comments")
			{
				comments.GET("", d.Authorize("comments", "get", d.adapterSlave), d.getComments)
				comments.GET("/:id", d.Authorize("comment", "get", d.adapterSlave), d.getComment)
				comments.POST("", d.Authorize("comment", "post", d.adapterSlave), d.createComment)
				comments.PATCH("", d.Authorize("comment", "patch", d.adapterSlave), d.updateComment)
				comments.DELETE("/:id", d.Authorize("comment", "delete", d.adapterSlave), d.deleteComment)
			}

			specifications := version1.Group("/specifications")
			{
				specifications.GET("", d.Authorize("specifications", "get", d.adapterSlave), d.getSpecifications)
				specifications.GET("/:id", d.Authorize("specification", "get", d.adapterSlave), d.getSpecification)
				specifications.POST("", d.Authorize("specification", "post", d.adapterSlave), d.createSpecification)
				specifications.PATCH("", d.Authorize("specification", "patch", d.adapterSlave), d.updateSpecification)
				specifications.DELETE("/:id", d.Authorize("specification", "delete", d.adapterSlave), d.deleteSpecification)
			}

			favorites := version1.Group("/favorites")
			{
				favorites.POST("", d.Authorize("favorite", "post", d.adapterSlave), d.createFavorite)
				favorites.DELETE("/:item_id", d.Authorize("favorite", "delete", d.adapterSlave), d.deleteFavorite)
			}

			rules := version1.Group("/rules")
			{
				rules.GET("", d.Authorize("rules", "get", d.adapterSlave), d.getRules)
				rules.GET("/:id", d.Authorize("rule", "get", d.adapterSlave), d.getRule)
				rules.POST("", d.Authorize("rule", "post", d.adapterSlave), d.createRule)
				rules.PATCH("", d.Authorize("rule", "patch", d.adapterSlave), d.updateRule)
				rules.DELETE("/:id", d.Authorize("rule", "delete", d.adapterSlave), d.deleteRule)
			}

			messages := version1.Group("/messages")
			{
				messages.GET("", d.Authorize("messages", "get", d.adapterSlave), d.getMessages)
				messages.GET("/:id", d.Authorize("message", "get", d.adapterSlave), d.getMessage)
				messages.POST("", d.Authorize("message", "post", d.adapterSlave), d.createMessage)
				messages.PATCH("", d.Authorize("message", "patch", d.adapterSlave), d.updateMessage)
				messages.DELETE("/:id", d.Authorize("message", "delete", d.adapterSlave), d.deleteMessage)
			}
		}
	}

	return router
}
