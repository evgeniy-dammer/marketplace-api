package model

// Comment model.
type Comment struct {
	ID             string  `json:"id" db:"id"`
	UserID         string  `json:"user" db:"user_created" binding:"required"`
	ItemID         string  `json:"item" db:"item_id" binding:"required"`
	Content        string  `json:"content" db:"content" binding:"required"`
	Status         int     `json:"status" db:"status_id" binding:"required"`
	Rating         float32 `json:"rating" db:"rating" binding:"required"`
	OrganizationID string  `json:"organization" db:"organization_id" binding:"required"`
	CreatedAt      string  `json:"created,omitempty" db:"created_at"`
}

// UpdateCommentInput is an input data for updating comment.
type UpdateCommentInput struct {
	ID             *string  `json:"id"`
	UserID         *string  `json:"user"`
	ItemID         *string  `json:"item"`
	Content        *string  `json:"content"`
	Status         *int     `json:"status"`
	Rating         *float32 `json:"rating"`
	OrganizationID *string  `json:"organization"`
}

// Validate checks if update input is nil.
func (i UpdateCommentInput) Validate() error {
	if i.ID == nil &&
		i.UserID == nil &&
		i.ItemID == nil &&
		i.Content == nil &&
		i.Status == nil &&
		i.OrganizationID == nil {
		return errStructHasNoValues
	}

	return nil
}
