package repository

import (
	"fmt"
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

// ItemPostgresql repository
type ItemPostgresql struct {
	db *sqlx.DB
}

// NewItemPostgresql is a constructor for ItemPostgresql
func NewItemPostgresql(db *sqlx.DB) *ItemPostgresql {
	return &ItemPostgresql{db: db}
}

// GetAll selects all items from database
func (r *ItemPostgresql) GetAll(userId string, organizationId string) ([]model.Item, error) {
	var items []model.Item

	query := fmt.Sprintf("SELECT id, name, price, category_id, organisation_id FROM %s WHERE organisation_id = $1 ",
		itemTable)
	err := r.db.Select(&items, query, organizationId)

	return items, err
}

// GetOne select item by id from database
func (r *ItemPostgresql) GetOne(userId string, organizationId string, itemId string) (model.Item, error) {
	var item model.Item

	query := fmt.Sprintf("SELECT id, name, price, category_id, organisation_id FROM %s WHERE organisation_id = $1 AND id = $2 ",
		itemTable)
	err := r.db.Get(&item, query, organizationId, itemId)

	return item, err
}

// Create insert item into database
func (r *ItemPostgresql) Create(userId string, organizationId string, item model.Item) (string, error) {
	var id string

	query := fmt.Sprintf("INSERT INTO %s (name, price, category_id, organisation_id) VALUES ($1, $2, $3, $4) RETURNING id",
		itemTable)
	row := r.db.QueryRow(query, item.Name, item.Price, item.CategoryId, organizationId)
	err := row.Scan(&id)

	return id, err
}

// Update updates item by id in database
func (r *ItemPostgresql) Update(userId string, organizationId string, itemId string, input model.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, *input.Price)
		argId++
	}

	if input.CategoryId != nil {
		setValues = append(setValues, fmt.Sprintf("category_id=$%d", argId))
		args = append(args, *input.CategoryId)
		argId++
	}

	if input.OrganisationId != nil {
		setValues = append(setValues, fmt.Sprintf("organisation_id=$%d", argId))
		args = append(args, *input.OrganisationId)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE organisation_id = '%s' AND id = '%s'",
		itemTable, setQuery, organizationId, itemId)
	_, err := r.db.Exec(query, args...)

	return err
}

// Delete deletes item by id from database
func (r *ItemPostgresql) Delete(userId string, organizationId string, itemId string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE organisation_id = $1 AND id = $2", itemTable)
	_, err := r.db.Exec(query, organizationId, itemId)

	return err
}
