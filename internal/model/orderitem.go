package model

// OrderItem model.
type OrderItem struct {
	ID         string  `json:"id"`
	OrderID    string  `json:"orderid"`
	Quantity   int     `json:"quantity"`
	ItemID     string  `json:"itemid"`
	UnitPrice  float32 `json:"unitprise"`
	TotalPrice float32 `json:"totalprice"`
}
