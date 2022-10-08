package controller

import (
	"github.com/evgeniy-dammer/emenu-api/pkg/items"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(app *fiber.App, db *pgx.Conn) {
	itemRoutes := app.Group("/api/item")
	itemRoutes.Get("/", items.GetAllItems)
	//routes.Get("/", h.GetUsers)
	//routes.Get("/:id", h.GetUser)
	//routes.Put("/:id", h.UpdateUser)
	//routes.Delete("/:id", h.DeleteUser)
}
