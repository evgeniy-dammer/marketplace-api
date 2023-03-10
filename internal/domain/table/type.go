package table

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// Table entity.
type Table struct {
	// Table ID
	ID string `json:"id" db:"id"`
	// Table name
	Name string `json:"name" db:"name" binding:"required"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
}

// CreateTableInput entity.
type CreateTableInput struct {
	// Table name
	Name string `json:"name" db:"name" binding:"required"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
}

// UpdateTableInput is an input data for updating table entity.
type UpdateTableInput struct {
	// Table ID
	ID *string `json:"id"`
	// Table name
	Name *string `json:"name"`
	// Organization ID
	OrganizationID *string `json:"organisation"`
}

// Validate checks if update input is nil.
func (i UpdateTableInput) Validate() error {
	if i.ID == nil && i.Name == nil && i.OrganizationID == nil {
		return ErrStructHasNoValues
	}

	return nil
}
