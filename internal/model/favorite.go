package model

// Favorite model.
type Favorite struct {
	ItemID string `json:"item" db:"item_id" binding:"required"`
}
