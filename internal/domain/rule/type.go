package rule

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// Rule entity
type Rule struct {
	// Rule ID
	ID string `json:"id" db:"id"`
	// Type
	Ptype string `json:"ptype" db:"ptype" binding:"required"`
	// Role name
	V0 string `json:"v0" db:"v0" binding:"required"`
	// Recourse name
	V1 string `json:"v1" db:"v1" binding:"required"`
	// REST verb
	V2 string `json:"v2" db:"v2" binding:"required"`
	// Permission status
	V3 string `json:"v3" db:"v3" binding:"required"`
	// Empty
	V4 string `json:"v4" db:"v4"`
	// Empty
	V5 string `json:"v5" db:"v5"`
}

// CreateRuleInput entity
type CreateRuleInput struct {
	// Type
	Ptype string `json:"ptype" db:"ptype" binding:"required"`
	// Role name
	V0 string `json:"v0" db:"v0" binding:"required"`
	// Recourse name
	V1 string `json:"v1" db:"v1" binding:"required"`
	// REST verb
	V2 string `json:"v2" db:"v2" binding:"required"`
	// Permission status
	V3 string `json:"v3" db:"v3" binding:"required"`
	// Empty
	V4 string `json:"v4" db:"v4"`
	// Empty
	V5 string `json:"v5" db:"v5"`
}

// UpdateRuleInput is an input data for updating rule entity.
type UpdateRuleInput struct {
	// Rule ID
	ID *string `json:"id"`
	// Type
	Ptype *string `json:"ptype"`
	// Role name
	V0 *string `json:"v0"`
	// Recourse name
	V1 *string `json:"v1"`
	// REST verb
	V2 *string `json:"v2"`
	// Permission status
	V3 *string `json:"v3"`
	// Empty
	V4 *string `json:"v4"`
	// Empty
	V5 *string `json:"v5"`
}

// Validate checks if update input is nil.
func (i UpdateRuleInput) Validate() error {
	if i.ID == nil &&
		i.Ptype == nil &&
		i.V0 == nil &&
		i.V1 == nil &&
		i.V2 == nil &&
		i.V3 == nil {
		return ErrStructHasNoValues
	}

	return nil
}
