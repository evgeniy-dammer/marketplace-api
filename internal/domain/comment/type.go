package comment

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// Comment entity.
type Comment struct {
	// Comment ID
	ID string `json:"id" db:"id"`
	// User ID
	UserID string `json:"user" db:"user_created" binding:"required"`
	// Item ID
	ItemID string `json:"item" db:"item_id" binding:"required"`
	// Comments content message
	Content string `json:"content" db:"content" binding:"required"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Created datetime
	CreatedAt string `json:"created,omitempty" db:"created_at"`
	// Comments status ID
	Status int `json:"status" db:"status_id" binding:"required"`
	// Comments rating value
	Rating float32 `json:"rating" db:"rating" binding:"required"`
}

// CreateCommentInput entity.
type CreateCommentInput struct {
	// User ID
	UserID string `json:"user" db:"user_created" binding:"required"`
	// Item ID
	ItemID string `json:"item" db:"item_id" binding:"required"`
	// Comments content message
	Content string `json:"content" db:"content" binding:"required"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Comments status ID
	Status int `json:"status" db:"status_id" binding:"required"`
	// Comments rating value
	Rating float32 `json:"rating" db:"rating" binding:"required"`
}

// UpdateCommentInput is an input data for updating comment entity.
type UpdateCommentInput struct {
	// Comment ID
	ID *string `json:"id"`
	// User ID
	UserID *string `json:"user"`
	// Item ID
	ItemID *string `json:"item"`
	// Comments content message
	Content *string `json:"content"`
	// Organization ID
	OrganizationID *string `json:"organization"`
	// Comments status ID
	Status *int `json:"status"`
	// Comments rating value
	Rating *float32 `json:"rating"`
}

// Validate checks if update input is nil.
func (i UpdateCommentInput) Validate() error {
	if i.ID == nil &&
		i.UserID == nil &&
		i.ItemID == nil &&
		i.Content == nil &&
		i.Status == nil &&
		i.OrganizationID == nil {
		return ErrStructHasNoValues
	}

	return nil
}
