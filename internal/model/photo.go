package model

// Photo model.
type Photo struct {
	ID        string `json:"id"`
	BigPath   string `json:"bigpath"`
	ThumbPath string `json:"thumbpath"`
}
