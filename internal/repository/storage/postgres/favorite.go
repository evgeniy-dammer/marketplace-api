package postgres

import (
	"fmt"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/favorite"
	"github.com/pkg/errors"
)

// FavoriteCreate insert favorite into database.
func (r *Repository) FavoriteCreate(userID string, favorite favorite.Favorite) error {
	var favoriteID string

	query := fmt.Sprintf("INSERT INTO %s (user_id, item_id) VALUES ($1, $2) RETURNING id", favoriteTable)

	row := r.db.QueryRow(query, userID, favorite.ItemID)
	err := row.Scan(&favoriteID)

	return errors.Wrap(err, "favorite create query error")
}

// FavoriteDelete deletes favorite by userID and itemID from database.
func (r *Repository) FavoriteDelete(userID string, itemID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND item_id = $2", favoriteTable)

	_, err := r.db.Exec(query, userID, itemID)

	return errors.Wrap(err, "favorite delete query error")
}
