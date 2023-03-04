package image

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/image"
	"github.com/pkg/errors"
)

// ImageGetAll returns all images from the system.
func (s *UseCase) ImageGetAll(userID string, organizationID string) ([]image.Image, error) {
	images, err := s.adapterStorage.ImageGetAll(userID, organizationID)

	return images, errors.Wrap(err, "images select error")
}

// ImageGetOne returns image by id from the system.
func (s *UseCase) ImageGetOne(userID string, organizationID string, imageID string) (image.Image, error) {
	imageSingle, err := s.adapterStorage.ImageGetOne(userID, organizationID, imageID)

	return imageSingle, errors.Wrap(err, "image select error")
}

// ImageCreate inserts image into system.
func (s *UseCase) ImageCreate(userID string, image image.Image) (string, error) {
	imageID, err := s.adapterStorage.ImageCreate(userID, image)

	return imageID, errors.Wrap(err, "image create error")
}

// ImageUpdate updates image by id in the system.
func (s *UseCase) ImageUpdate(userID string, input image.UpdateImageInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.ImageUpdate(userID, input)

	return errors.Wrap(err, "image update error")
}

// ImageDelete deletes image by id from the system.
func (s *UseCase) ImageDelete(userID string, organizationID string, imageID string) error {
	err := s.adapterStorage.ImageDelete(userID, organizationID, imageID)

	return errors.Wrap(err, "image delete error")
}
