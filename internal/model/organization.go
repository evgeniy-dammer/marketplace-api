package model

// Organization model.
type Organization struct {
	ID      string `json:"id" db:"id"`
	Name    string `json:"name" db:"name" binding:"required"`
	UserID  string `json:"userid" db:"user_id" binding:"required"`
	Address string `json:"address" db:"address" binding:"required"`
	Phone   string `json:"phone" db:"phone" binding:"required"`
}

// UpdateOrganizationInput is an input data for updating organization.
type UpdateOrganizationInput struct {
	ID      *string `json:"id"`
	Name    *string `json:"name"`
	Address *string `json:"address"`
	Phone   *string `json:"phone"`
}

// Validate checks if update input is nil.
func (i UpdateOrganizationInput) Validate() error {
	if i.ID == nil && i.Name == nil {
		return errStructHasNoValues
	}

	return nil
}
