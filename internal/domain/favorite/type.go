package favorite

// Favorite entity.
type Favorite struct {
	ItemID string `json:"item" db:"item_id" binding:"required"`
}
