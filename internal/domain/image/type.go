package image

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// Image entity.
type Image struct {
	ID             string `json:"id"`
	ObjectID       string `json:"object" db:"object_id" binding:"required"`
	Origin         string `json:"origin" db:"origin" binding:"required"`
	Middle         string `json:"middle" db:"middle" binding:"required"`
	Small          string `json:"small" db:"small" binding:"required"`
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	ObjectType     int    `json:"type" db:"type" binding:"required"`
	IsMain         bool   `json:"main" db:"is_main" binding:"required"`
}

// UpdateImageInput is an input data for updating image entity.
type UpdateImageInput struct {
	ID             *string `json:"id"`
	ObjectID       *string `json:"object"`
	Origin         *string `json:"origin"`
	Middle         *string `json:"middle"`
	Small          *string `json:"small"`
	OrganizationID *string `json:"organization"`
	ObjectType     *int    `json:"type"`
	IsMain         *bool   `json:"main"`
}

// Validate checks if update input is nil.
func (i UpdateImageInput) Validate() error {
	if i.ID == nil && i.ObjectID == nil && i.ObjectType == nil && i.OrganizationID == nil {
		return ErrStructHasNoValues
	}

	return nil
}
