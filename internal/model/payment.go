package model

// Payment model.
type Payment struct {
	ID             string  `json:"id"`
	DateTime       string  `json:"datetime"`
	Sum            float32 `json:"sum"`
	OrderID        string  `json:"orderid"`
	DepartmentID   string  `json:"departmentid"`
	OrganisationID string  `json:"organisationid"`
	Status         int     `json:"status"`
}
