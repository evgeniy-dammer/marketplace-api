package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/pkg/errors"
)

// ImageService is an image service.
type ImageService struct {
	repo repository.Image
}

// NewImageService is a constructor for ImageService.
func NewImageService(repo repository.Image) *ImageService {
	return &ImageService{repo: repo}
}

// GetAll returns all images from the system.
func (s *ImageService) GetAll(userID string, organizationID string) ([]model.Image, error) {
	images, err := s.repo.GetAll(userID, organizationID)

	return images, errors.Wrap(err, "images select error")
}

// GetOne returns image by id from the system.
func (s *ImageService) GetOne(userID string, organizationID string, imageID string) (model.Image, error) {
	image, err := s.repo.GetOne(userID, organizationID, imageID)

	return image, errors.Wrap(err, "image select error")
}

// Create inserts image into system.
func (s *ImageService) Create(userID string, image model.Image) (string, error) {
	imageID, err := s.repo.Create(userID, image)

	return imageID, errors.Wrap(err, "image create error")
}

// Update updates image by id in the system.
func (s *ImageService) Update(userID string, input model.UpdateImageInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.repo.Update(userID, input)

	return errors.Wrap(err, "image update error")
}

// Delete deletes image by id from the system.
func (s *ImageService) Delete(userID string, organizationID string, imageID string) error {
	err := s.repo.Delete(userID, organizationID, imageID)

	return errors.Wrap(err, "image delete error")
}
