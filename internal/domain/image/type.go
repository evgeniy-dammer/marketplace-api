package image

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// Image entity.
type Image struct {
	// Image ID
	ID string `json:"id"`
	// Object ID
	ObjectID string `json:"object" db:"object_id" binding:"required"`
	// Origin file name
	Origin string `json:"origin" db:"origin" binding:"required"`
	// Middle-sized file name
	Middle string `json:"middle" db:"middle" binding:"required"`
	// Small-sized file name
	Small string `json:"small" db:"small" binding:"required"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Object type
	ObjectType int `json:"type" db:"type" binding:"required"`
	// Is main image
	IsMain bool `json:"main" db:"is_main" binding:"required"`
}

// CreateImageInput entity.
type CreateImageInput struct {
	// Object ID
	ObjectID string `json:"object" db:"object_id" binding:"required"`
	// Origin file name
	Origin string `json:"origin" db:"origin" binding:"required"`
	// Middle-sized file name
	Middle string `json:"middle" db:"middle" binding:"required"`
	// Small-sized file name
	Small string `json:"small" db:"small" binding:"required"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Object type
	ObjectType int `json:"type" db:"type" binding:"required"`
	// Is main image
	IsMain bool `json:"main" db:"is_main" binding:"required"`
}

// UpdateImageInput is an input data for updating image entity.
type UpdateImageInput struct {
	// Image ID
	ID *string `json:"id"`
	// Object ID
	ObjectID *string `json:"object"`
	// Origin file name
	Origin *string `json:"origin"`
	// Middle-sized file name
	Middle *string `json:"middle"`
	// Small-sized file name
	Small *string `json:"small"`
	// Organization ID
	OrganizationID *string `json:"organization"`
	// Object type
	ObjectType *int `json:"type"`
	// Is main image
	IsMain *bool `json:"main"`
}

// Validate checks if update input is nil.
func (i UpdateImageInput) Validate() error {
	if i.ID == nil && i.ObjectID == nil && i.ObjectType == nil && i.OrganizationID == nil {
		return ErrStructHasNoValues
	}

	return nil
}
