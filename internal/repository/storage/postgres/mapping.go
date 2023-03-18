package postgres

import "github.com/evgeniy-dammer/marketplace-api/pkg/columncode"

var mappingSortUser = map[columncode.ColumnCode]string{
	"id":         "id",
	"phone":      "phone",
	"first_name": "first_name",
	"last_name":  "last_name",
	"status_id":  "status_id",
	"created_at": "created_at",
}

var mappingSortRole = map[columncode.ColumnCode]string{
	"id":   "id",
	"name": "name",
}
