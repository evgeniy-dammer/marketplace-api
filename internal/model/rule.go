package model

type Rule struct {
	ID    string `json:"id" db:"id"`
	Ptype string `json:"ptype" db:"ptype" binding:"required"`
	V0    string `json:"v0" db:"v0" binding:"required"`
	V1    string `json:"v1" db:"v1" binding:"required"`
	V2    string `json:"v2" db:"v2" binding:"required"`
	V3    string `json:"v3" db:"v3" binding:"required"`
	V4    string `json:"v4" db:"v4"`
	V5    string `json:"v5" db:"v5"`
}

// UpdateRuleInput is an input data for updating rule.
type UpdateRuleInput struct {
	ID    *string `json:"id"`
	Ptype *string `json:"ptype"`
	V0    *string `json:"v0"`
	V1    *string `json:"v1"`
	V2    *string `json:"v2"`
	V3    *string `json:"v3"`
	V4    *string `json:"v4"`
	V5    *string `json:"v5"`
}

// Validate checks if update input is nil.
func (i UpdateRuleInput) Validate() error {
	if i.ID == nil &&
		i.Ptype == nil &&
		i.V0 == nil &&
		i.V1 == nil &&
		i.V2 == nil &&
		i.V3 == nil {
		return errStructHasNoValues
	}

	return nil
}
