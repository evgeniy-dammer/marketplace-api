package model

// Department model
type Department struct {
	Id             string `json:"id" db:"id"`
	Name           string `json:"name" db:"name" binding:"required"`
	Address        string `json:"address" db:"address" binding:"required"`
	Phone          string `json:"phone" db:"phone" binding:"required"`
	OrganisationId string `json:"organisation" db:"organisation_id" binding:"required"`
}
