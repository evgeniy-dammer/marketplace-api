package model

// Order model.
type Order struct {
	ID             string      `json:"id"`
	DateTime       string      `json:"datetime"`
	TableID        string      `json:"tableid"`
	Status         int         `json:"status"`
	OrganisationID string      `json:"organisationid"`
	DepartmentID   string      `json:"departmentid"`
	TotalSum       float32     `json:"totalsum"`
	Items          []OrderItem `json:"orderitems"`
}
