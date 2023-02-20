package model

// Table model.
type Table struct {
	ID             string `json:"id" db:"id"`
	Name           string `json:"name" db:"name" binding:"required"`
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
}
