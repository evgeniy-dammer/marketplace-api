package favorite

// ListFavorite
//
//easyjson:json
type ListFavorite []Favorite

// Favorite entity.
//
//easyjson:json
type Favorite struct {
	// Item ID
	ItemID string `json:"item" db:"item_id" binding:"required"`
}
