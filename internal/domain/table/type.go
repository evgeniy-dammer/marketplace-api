package table

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// Table entity.
type Table struct {
	ID             string `json:"id" db:"id"`
	Name           string `json:"name" db:"name" binding:"required"`
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
}

// UpdateTableInput is an input data for updating table entity.
type UpdateTableInput struct {
	ID             *string `json:"id"`
	Name           *string `json:"name"`
	OrganizationID *string `json:"organisation"`
}

// Validate checks if update input is nil.
func (i UpdateTableInput) Validate() error {
	if i.ID == nil && i.Name == nil && i.OrganizationID == nil {
		return ErrStructHasNoValues
	}

	return nil
}
