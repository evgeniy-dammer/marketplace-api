package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/pkg/errors"
)

// FavoriteService is an organization service.
type FavoriteService struct {
	repo repository.Favorite
}

// NewFavoriteService is a constructor for FavoriteService.
func NewFavoriteService(repo repository.Favorite) *FavoriteService {
	return &FavoriteService{repo: repo}
}

// Create inserts favorite into system.
func (s *FavoriteService) Create(userID string, favorite model.Favorite) error {
	err := s.repo.Create(userID, favorite)

	return errors.Wrap(err, "favorite create error")
}

// Delete deletes favorite by id from the system.
func (s *FavoriteService) Delete(userID string, itemID string) error {
	err := s.repo.Delete(userID, itemID)

	return errors.Wrap(err, "favorite delete error")
}
