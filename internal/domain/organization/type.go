package organization

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// ListOrganization
//
//easyjson:json
type ListOrganization []Organization

// Organization entity.
//
//easyjson:json
type Organization struct {
	// Organization ID
	ID string `json:"id" db:"id"`
	// Organization name
	Name string `json:"name" db:"name" binding:"required"`
	// User ID
	UserID string `json:"userid" db:"user_id" binding:"required"`
	// Organization address
	Address string `json:"address" db:"address" binding:"required"`
	// Organization phone number
	Phone string `json:"phone" db:"phone" binding:"required"`
}

// CreateOrganizationInput entity.
//
//easyjson:json
type CreateOrganizationInput struct {
	// Organization name
	Name string `json:"name" db:"name" binding:"required"`
	// User ID
	UserID string `json:"userid" db:"user_id" binding:"required"`
	// Organization address
	Address string `json:"address" db:"address" binding:"required"`
	// Organization phone number
	Phone string `json:"phone" db:"phone" binding:"required"`
}

// UpdateOrganizationInput is an input data for updating organization entity.
//
//easyjson:json
type UpdateOrganizationInput struct {
	// Organization ID
	ID *string `json:"id"`
	// Organization name
	Name *string `json:"name"`
	// Organization address
	Address *string `json:"address"`
	// Organization phone number
	Phone *string `json:"phone"`
}

// Validate checks if update input is nil.
func (i UpdateOrganizationInput) Validate() error {
	if i.ID == nil && i.Name == nil {
		return ErrStructHasNoValues
	}

	return nil
}
