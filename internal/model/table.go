package model

// Table model
type Table struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	OrganisationId string `json:"organisationid"`
	DepartmentId   string `json:"departmentid"`
}
