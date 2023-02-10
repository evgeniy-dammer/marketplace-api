package model

// Item model
type Item struct {
	Id             string  `json:"id" db:"id"`
	Name           string  `json:"name" db:"name" binding:"required"`
	Price          float32 `json:"price" db:"price"`
	CategoryId     string  `json:"category" db:"category_id" binding:"required"`
	OrganisationId string  `json:"organisation" db:"organisation_id" binding:"required"`
}
