package models

//Item model
type Item struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	Price          float32 `json:"price"`
	CategoryId     string  `json:"categoryid"`
	OrganisationId string  `json:"organisationid"`
	PrepareTime    int     `json:"preparetime"`
}
