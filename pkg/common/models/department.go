package models

type Department struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	OrganisationId string `json:"organisationid"`
}
