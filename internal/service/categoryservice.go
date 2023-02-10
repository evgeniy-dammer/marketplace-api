package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
)

// CategoryService is a category service
type CategoryService struct {
	repo repository.Category
}

// NewCategoryService is a constructor for CategoryService
func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetAll returns all categories from the system
func (s *CategoryService) GetAll(userId string, organizationId string) ([]model.Category, error) {
	return s.repo.GetAll(userId, organizationId)
}

// GetOne returns category by id from the system
func (s *CategoryService) GetOne(userId string, organizationId string, categoryId string) (model.Category, error) {
	return s.repo.GetOne(userId, organizationId, categoryId)
}

// Create inserts category into system
func (s *CategoryService) Create(userId string, organizationId string, category model.Category) (string, error) {
	return s.repo.Create(userId, organizationId, category)
}

// Update updates category by id in the system
func (s *CategoryService) Update(userId string, organizationId string, categoryId string, input model.UpdateCategoryInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, organizationId, categoryId, input)
}

// Delete deletes category by id from the system
func (s *CategoryService) Delete(userId string, organizationId string, categoryId string) error {
	return s.repo.Delete(userId, organizationId, categoryId)
}
