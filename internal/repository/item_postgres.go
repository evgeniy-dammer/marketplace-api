package repository

import (
	"fmt"
	"strings"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ItemPostgresql repository.
type ItemPostgresql struct {
	db *sqlx.DB
}

// NewItemPostgresql is a constructor for ItemPostgresql.
func NewItemPostgresql(db *sqlx.DB) *ItemPostgresql {
	return &ItemPostgresql{db: db}
}

// GetAll selects all items from database.
func (r *ItemPostgresql) GetAll(userID string, organizationID string) ([]model.Item, error) {
	var items []model.Item

	query := fmt.Sprintf("SELECT id, name, price, category_id, organisation_id FROM %s WHERE organisation_id = $1 ",
		itemTable)
	err := r.db.Select(&items, query, organizationID)

	return items, errors.Wrap(err, "items select query error")
}

// GetOne select item by id from database.
func (r *ItemPostgresql) GetOne(userID string, organizationID string, itemID string) (model.Item, error) {
	var item model.Item

	query := fmt.Sprintf(
		"SELECT id, name, price, category_id, organisation_id FROM %s WHERE organisation_id = $1 AND id = $2 ",
		itemTable,
	)
	err := r.db.Get(&item, query, organizationID, itemID)

	return item, errors.Wrap(err, "item select query error")
}

// Create insert item into database.
func (r *ItemPostgresql) Create(userID string, organizationID string, item model.Item) (string, error) {
	var itemID string

	query := fmt.Sprintf("INSERT INTO %s (name, price, category_id, organisation_id) VALUES ($1, $2, $3, $4) RETURNING id",
		itemTable)
	row := r.db.QueryRow(query, item.Name, item.Price, item.CategoryID, organizationID)
	err := row.Scan(&itemID)

	return itemID, errors.Wrap(err, "item create query error")
}

// Update updates item by id in database.
func (r *ItemPostgresql) Update(userID string, organizationID string, itemID string, input model.UpdateItemInput) error { //nolint:lll
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argID))
		args = append(args, *input.Name)
		argID++
	}

	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argID))
		args = append(args, *input.Price)
		argID++
	}

	if input.CategoryID != nil {
		setValues = append(setValues, fmt.Sprintf("category_id=$%d", argID))
		args = append(args, *input.CategoryID)
		argID++
	}

	if input.OrganisationID != nil {
		setValues = append(setValues, fmt.Sprintf("organisation_id=$%d", argID))
		args = append(args, *input.OrganisationID)
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE organisation_id = '%s' AND id = '%s'",
		itemTable, setQuery, organizationID, itemID)
	_, err := r.db.Exec(query, args...)

	return errors.Wrap(err, "item update query error")
}

// Delete deletes item by id from database.
func (r *ItemPostgresql) Delete(userID string, organizationID string, itemID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE organisation_id = $1 AND id = $2", itemTable)
	_, err := r.db.Exec(query, organizationID, itemID)

	return errors.Wrap(err, "item delete query error")
}
