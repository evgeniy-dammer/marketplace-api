package model

// Category model.
type Category struct {
	ID             string `json:"id" db:"id"`
	Name           string `json:"name" db:"name" binding:"required"`
	Parent         string `json:"parent" db:"parent_id"`
	Level          int    `json:"level" db:"level"`
	OrganisationID string `json:"organisation" db:"organisation_id" binding:"required"`
}
