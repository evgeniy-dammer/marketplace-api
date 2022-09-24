package models

//Order model
type Order struct {
	Id             string      `json:"id"`
	DateTime       string      `json:"datetime"`
	TableId        string      `json:"tableid"`
	Status         int         `json:"status"`
	OrganisationId string      `json:"organisationid"`
	DepartmentId   string      `json:"departmentid"`
	TotalSum       float32     `json:"totalsum"`
	Items          []OrderItem `json:"orderitems"`
}
