package favorite

// Favorite entity.
type Favorite struct {
	// Item ID
	ItemID string `json:"item" db:"item_id" binding:"required"`
}
