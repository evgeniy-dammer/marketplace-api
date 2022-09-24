package models

type Category struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Parent         string `json:"parent"`
	Level          int    `json:"level"`
	OrganisationId string `json:"organisationid"`
}
