package repository

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
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
	/*
		query := fmt.Sprintf("SELECT id, name, price, category_id, organisation_id "+
			"FROM %s WHERE organisation_id = $1 ",
			itemTable)

		err := r.db.Select(&items, query, organizationId)
	*/
	return items, nil
}

// GetOne select item by id from database
func (r *ItemPostgresql) GetOne(userId string, organizationId string, itemId string) (model.Item, error) {
	var item model.Item
	/*
		query := fmt.Sprintf(
			"SELECT id, name, price, category_id, organisation_id FROM %s WHERE organisation_id = $1 AND id = $2 ",
			itemTable)
		err := r.db.Get(&item, query, organizationId, itemId)
	*/
	return item, nil
}

// Create insert item into database
func (r *ItemPostgresql) Create(userId string, organizationId string, item model.Item) (string, error) {
	/*var id string

	createUserQuery := fmt.Sprintf(
		"INSERT INTO %s (name, user_id, address, phone) VALUES ($1, $2, $3, $4) RETURNING id",
		organisationTable)

	row := r.db.QueryRow(createUserQuery, item.Name, userId, item.Address, item.Phone)

	err := row.Scan(&id)
	*/
	return "", nil
}

// Update updates item by id in database
func (r *ItemPostgresql) Update(userId string, organizationId string, itemId string, input model.UpdateItemInput) error {
	/*setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Address != nil {
		setValues = append(setValues, fmt.Sprintf("address=$%d", argId))
		args = append(args, *input.Address)
		argId++
	}

	if input.Phone != nil {
		setValues = append(setValues, fmt.Sprintf("phone=$%d", argId))
		args = append(args, *input.Phone)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = '%s' AND user_id = '%s'",
		organisationTable, setQuery, organizationId, userId)

	_, err := r.db.Exec(query, args...)
	*/
	return nil
}

// Delete deletes item by id from database
func (r *ItemPostgresql) Delete(userId string, organizationId string, itemId string) error {
	/*query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", organisationTable)
	_, err := r.db.Exec(query, organizationId, userId)
	*/
	return nil
}
