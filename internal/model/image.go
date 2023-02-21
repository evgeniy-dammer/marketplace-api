package model

// Image model.
type Image struct {
	ID     string `json:"id"`
	Origin string `json:"origin"`
	Middle string `json:"middle"`
	Small  string `json:"small"`
}
