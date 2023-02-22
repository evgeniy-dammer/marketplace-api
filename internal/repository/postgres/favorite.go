package postgres

import (
	"fmt"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// FavoritePostgresql repository.
type FavoritePostgresql struct {
	db *sqlx.DB
}

// NewFavoritePostgresql is a constructor for FavoritePostgresql.
func NewFavoritePostgresql(db *sqlx.DB) *FavoritePostgresql {
	return &FavoritePostgresql{db: db}
}

// Create insert favorite into database.
func (r *FavoritePostgresql) Create(userID string, favorite model.Favorite) error {
	var favoriteID string

	query := fmt.Sprintf("INSERT INTO %s (user_id, item_id) VALUES ($1, $2) RETURNING id", favoriteTable)

	row := r.db.QueryRow(query, userID, favorite.ItemID)
	err := row.Scan(&favoriteID)

	return errors.Wrap(err, "favorite create query error")
}

// Delete deletes favorite by userID and itemID from database.
func (r *FavoritePostgresql) Delete(userID string, itemID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND item_id = $2", favoriteTable)

	_, err := r.db.Exec(query, userID, itemID)

	return errors.Wrap(err, "favorite delete query error")
}
