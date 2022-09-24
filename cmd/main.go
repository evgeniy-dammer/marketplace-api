package main

import (
	"log"

	"github.com/evgeniy-dammer/emenu-api/pkg/common/config"
	"github.com/evgeniy-dammer/emenu-api/pkg/common/db"
)

func main() {
	configuration, err := config.LoadConfiguration()

	if err != nil {
		log.Fatalln("Configuration faild", err)
	}

	_ = db.Connect(&configuration)
}
