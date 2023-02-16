package model

// Table model.
type Table struct {
	ID             string `json:"id" db:"id"`
	Name           string `json:"name" db:"name" binding:"required"`
	OrganisationID string `json:"organisation" db:"organisation_id" binding:"required"`
}
