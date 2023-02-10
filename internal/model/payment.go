package model

// Payment model
type Payment struct {
	Id             string  `json:"id"`
	DateTime       string  `json:"datetime"`
	Sum            float32 `json:"sum"`
	OrderId        string  `json:"orderid"`
	DepartmentId   string  `json:"departmentid"`
	OrganisationId string  `json:"organisationid"`
	Status         int     `json:"status"`
}
