package main

import (
	"log"

	"github.com/evgeniy-dammer/emenu-api/pkg/common/config"
	"github.com/evgeniy-dammer/emenu-api/pkg/common/controller"
	"github.com/evgeniy-dammer/emenu-api/pkg/common/db"
	"github.com/gofiber/fiber/v2"
)

func main() {
	configuration, err := config.LoadConfiguration()

	if err != nil {
		log.Fatalln("Configuration faild", err)
	}

	db.Connect(&configuration)
	app := fiber.New()

	controller.RegisterRoutes(app, db.DB)

	app.Listen(":" + configuration.SvPort)
}
