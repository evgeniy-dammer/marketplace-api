package model

// Table model.
type Table struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganisationID string `json:"organisationid"`
	DepartmentID   string `json:"departmentid"`
}
