package model

// Category model
type Category struct {
	Id             string `json:"id" db:"id"`
	Name           string `json:"name" db:"name" binding:"required"`
	Parent         string `json:"parent" db:"parent_id"`
	Level          int    `json:"level" db:"level"`
	OrganisationId string `json:"organisation" db:"organisation_id" binding:"required"`
}
