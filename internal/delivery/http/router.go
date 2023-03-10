package http

import (
	"github.com/evgeniy-dammer/emenu-api/docs"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRoutes crete routes.
func (d *Delivery) InitRoutes(mode string) *gin.Engine {
	gin.SetMode(mode)

	router := gin.New()
	pprof.Register(router, "dev/pprof")
	router.Use(d.corsMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(ginzap.RecoveryWithZap(logger.Logger, true))
	router.RedirectTrailingSlash = false

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
				users.GET("/", d.Authorize("users", "get", d.adapter), d.getAllUsers)
				users.GET("/:id", d.Authorize("user", "get", d.adapter), d.getUser)
				users.POST("/", d.Authorize("user", "post", d.adapter), d.createUser)
				users.PATCH("/", d.Authorize("user", "patch", d.adapter), d.updateUser)
				users.DELETE("/:id", d.Authorize("user", "delete", d.adapter), d.deleteUser)
				users.GET("/roles", d.Authorize("roles", "get", d.adapter), d.getAllRoles)
			}

			organizations := version1.Group("/organizations")
			{
				organizations.GET("/", d.Authorize("organizations", "get", d.adapter), d.getOrganizations)
				organizations.GET("/:id", d.Authorize("organization", "get", d.adapter), d.getOrganization)
				organizations.POST("/", d.Authorize("organization", "post", d.adapter), d.createOrganization)
				organizations.PATCH("/", d.Authorize("organization", "patch", d.adapter), d.updateOrganization)
				organizations.DELETE("/:id", d.Authorize("organization", "delete", d.adapter), d.deleteOrganization)
			}

			categories := version1.Group("/categories")
			{
				categories.GET("/:org_id", d.Authorize("categories", "get", d.adapter), d.getCategories)
				categories.GET("/:org_id/:id", d.Authorize("category", "get", d.adapter), d.getCategory)
				categories.POST("/", d.Authorize("category", "post", d.adapter), d.createCategory)
				categories.PATCH("/", d.Authorize("category", "patch", d.adapter), d.updateCategory)
				categories.DELETE("/:org_id/:id", d.Authorize("category", "delete", d.adapter), d.deleteCategory)
			}

			items := version1.Group("/items")
			{
				items.GET("/:org_id", d.Authorize("items", "get", d.adapter), d.getItems)
				items.GET("/:org_id/:id", d.Authorize("item", "get", d.adapter), d.getItem)
				items.POST("/", d.Authorize("item", "post", d.adapter), d.createItem)
				items.PATCH("/", d.Authorize("item", "patch", d.adapter), d.updateItem)
				items.DELETE("/:org_id/:id", d.Authorize("item", "delete", d.adapter), d.deleteItem)
			}

			tables := version1.Group("/tables")
			{
				tables.GET("/:org_id", d.Authorize("tables", "get", d.adapter), d.getTables)
				tables.GET("/:org_id/:id", d.Authorize("table", "get", d.adapter), d.getTable)
				tables.POST("/", d.Authorize("table", "post", d.adapter), d.createTable)
				tables.PATCH("/", d.Authorize("table", "patch", d.adapter), d.updateTable)
				tables.DELETE("/:org_id/:id", d.Authorize("table", "delete", d.adapter), d.deleteTable)
			}

			orders := version1.Group("/orders")
			{
				orders.GET("/:org_id", d.Authorize("orders", "get", d.adapter), d.getOrders)
				orders.GET("/:org_id/:id", d.Authorize("order", "get", d.adapter), d.getOrder)
				orders.POST("/", d.Authorize("order", "post", d.adapter), d.createOrder)
				orders.PATCH("/", d.Authorize("order", "patch", d.adapter), d.updateOrder)
				orders.DELETE("/:org_id/:id", d.Authorize("order", "delete", d.adapter), d.deleteOrder)
			}

			images := version1.Group("/images")
			{
				images.GET("/:org_id", d.Authorize("images", "get", d.adapter), d.getImages)
				images.GET("/:org_id/:id", d.Authorize("image", "get", d.adapter), d.getImage)
				images.POST("/", d.Authorize("image", "post", d.adapter), d.createImage)
				images.PATCH("/", d.Authorize("image", "patch", d.adapter), d.updateImage)
				images.DELETE("/:org_id/:id", d.Authorize("image", "delete", d.adapter), d.deleteImage)
			}

			comments := version1.Group("/comments")
			{
				comments.GET("/:org_id", d.Authorize("comments", "get", d.adapter), d.getComments)
				comments.GET("/:org_id/:id", d.Authorize("comment", "get", d.adapter), d.getComment)
				comments.POST("/", d.Authorize("comment", "post", d.adapter), d.createComment)
				comments.PATCH("/", d.Authorize("comment", "patch", d.adapter), d.updateComment)
				comments.DELETE("/:org_id/:id", d.Authorize("comment", "delete", d.adapter), d.deleteComment)
			}

			specifications := version1.Group("/specifications")
			{
				specifications.GET("/:org_id", d.Authorize("specifications", "get", d.adapter), d.getSpecifications)
				specifications.GET("/:org_id/:id", d.Authorize("specification", "get", d.adapter), d.getSpecification)
				specifications.POST("/", d.Authorize("specification", "post", d.adapter), d.createSpecification)
				specifications.PATCH("/", d.Authorize("specification", "patch", d.adapter), d.updateSpecification)
				specifications.DELETE("/:org_id/:id", d.Authorize("specification", "delete", d.adapter), d.deleteSpecification)
			}

			favorites := version1.Group("/favorites")
			{
				favorites.POST("/", d.Authorize("favorite", "post", d.adapter), d.createFavorite)
				favorites.DELETE("/:item_id", d.Authorize("favorite", "delete", d.adapter), d.deleteFavorite)
			}

			rules := version1.Group("/rules")
			{
				rules.GET("/", d.Authorize("rules", "get", d.adapter), d.getRules)
				rules.GET("/:id", d.Authorize("rule", "get", d.adapter), d.getRule)
				rules.POST("/", d.Authorize("rule", "post", d.adapter), d.createRule)
				rules.PATCH("/", d.Authorize("rule", "patch", d.adapter), d.updateRule)
				rules.DELETE("/:id", d.Authorize("rule", "delete", d.adapter), d.deleteRule)
			}
		}
	}

	return router
}
