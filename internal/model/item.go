package model

// Item model
type Item struct {
	Id             string  `json:"id" db:"id"`
	Name           string  `json:"name" db:"name" binding:"required"`
	Price          float32 `json:"price" db:"price" binding:"required"`
	CategoryId     string  `json:"category" binding:"required"`
	OrganisationId string  `json:"organisation" db:"organisation_id" binding:"required"`
}
