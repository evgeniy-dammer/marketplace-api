package model

// Image model.
type Image struct {
	ID             string `json:"id"`
	ObjectID       string `json:"object" db:"object_id" binding:"required"`
	ObjectType     int    `json:"type" db:"type" binding:"required"`
	Origin         string `json:"origin" db:"origin" binding:"required"`
	Middle         string `json:"middle" db:"middle" binding:"required"`
	Small          string `json:"small" db:"small" binding:"required"`
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
}

// UpdateImageInput is an input data for updating image.
type UpdateImageInput struct {
	ID             *string `json:"id"`
	ObjectID       *string `json:"object"`
	ObjectType     *int    `json:"type"`
	Origin         *string `json:"origin"`
	Middle         *string `json:"middle"`
	Small          *string `json:"small"`
	OrganizationID *string `json:"organization"`
}

// Validate checks if update input is nil.
func (i UpdateImageInput) Validate() error {
	if i.ID == nil && i.ObjectID == nil && i.ObjectType == nil && i.OrganizationID == nil {
		return errStructHasNoValues
	}

	return nil
}
