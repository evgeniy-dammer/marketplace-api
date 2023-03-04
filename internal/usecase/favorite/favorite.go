package favorite

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/favorite"
	"github.com/pkg/errors"
)

// FavoriteCreate inserts favorite into system.
func (s *UseCase) FavoriteCreate(userID string, favorite favorite.Favorite) error {
	err := s.adapterStorage.FavoriteCreate(userID, favorite)

	return errors.Wrap(err, "favorite create error")
}

// FavoriteDelete deletes favorite by id from the system.
func (s *UseCase) FavoriteDelete(userID string, itemID string) error {
	err := s.adapterStorage.FavoriteDelete(userID, itemID)

	return errors.Wrap(err, "favorite delete error")
}
