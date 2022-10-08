package items

import (
	"context"
	"fmt"
	"os"

	"github.com/evgeniy-dammer/emenu-api/pkg/common/db"
	"github.com/evgeniy-dammer/emenu-api/pkg/common/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllItems(c *fiber.Ctx) error {
	var items []models.Item

	rows, err := db.DB.Query(context.Background(), "SELECT * FROM items")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Getting items failed: %v\n", err)
	}

	for rows.Next() {
		var item models.Item
		err = rows.Scan(&item.Id, &item.Name, &item.Price, &item.CategoryId, &item.OrganisationId, &item.PrepareTime)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Getting items failed: %v\n", err)
		}

		items = append(items, item)
	}

	return c.Status(fiber.StatusOK).JSON(&items)
}
