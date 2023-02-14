package model

// Department model.
type Department struct {
	ID             string `json:"id" db:"id"`
	Name           string `json:"name" db:"name" binding:"required"`
	Address        string `json:"address" db:"address" binding:"required"`
	Phone          string `json:"phone" db:"phone" binding:"required"`
	OrganisationID string `json:"organisation" db:"organisation_id" binding:"required"`
}
