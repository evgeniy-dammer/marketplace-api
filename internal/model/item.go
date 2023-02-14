package model

// Item model.
type Item struct {
	ID             string  `json:"id" db:"id"`
	Name           string  `json:"name" db:"name" binding:"required"`
	Price          float32 `json:"price" db:"price"`
	CategoryID     string  `json:"category" db:"category_id" binding:"required"`
	OrganisationID string  `json:"organisation" db:"organisation_id" binding:"required"`
}
