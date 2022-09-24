package models

type Photo struct {
	Id        string `json:"id"`
	BigPath   string `json:"bigpath"`
	ThumbPath string `json:"thumbpath"`
}
