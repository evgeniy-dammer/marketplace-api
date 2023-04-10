package message

import "github.com/pkg/errors"

var ErrStructHasNoValues = errors.New("update structure has no values")

//easyjson:json
type ListMessage []Message

// Message entity.
//
//easyjson:json
type Message struct {
	// Item ID
	ID string `json:"id" db:"id"`
	// Title
	Title string `json:"title" db:"title"`
	// Body
	Body string `json:"body" db:"body"`
	// User ID
	UserID string `json:"userid" db:"user_id"`
	// Is public
	IsPublic bool `json:"ispublic" db:"is_public" binding:"required"`
	// Is published
	IsPublished bool `json:"ispublished" db:"is_published" binding:"required"`
}

// CreateMessageInput entity
//
//easyjson:json
type CreateMessageInput struct {
	// Title
	Title string `json:"title" db:"title"`
	// Body
	Body string `json:"body" db:"body"`
	// User ID
	UserID string `json:"userid" db:"user_id"`
	// Is public
	IsPublic bool `json:"ispublic" db:"is_public" binding:"required"`
	// Is published
	IsPublished bool `json:"ispublished" db:"is_published" binding:"required"`
}

// UpdateMessageInput entity
//
//easyjson:json
type UpdateMessageInput struct {
	// Rule ID
	ID *string `json:"id"`
	// Title
	Title *string `json:"title"`
	// Body
	Body *string `json:"body"`
	// User ID
	UserID *string `json:"userid"`
	// Is public
	IsPublic *bool `json:"ispublic"`
	// Is published
	IsPublished *bool `json:"ispublished"`
}

// Validate checks if update input is nil.
func (i UpdateMessageInput) Validate() error {
	if i.ID == nil &&
		i.Title == nil &&
		i.Body == nil {
		return ErrStructHasNoValues
	}

	return nil
}
