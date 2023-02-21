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
