package models

type OrderItem struct {
	Id         string  `json:"id"`
	OrderId    string  `json:"orderid"`
	Quantity   int     `json:"quantity"`
	ItemId     string  `json:"itemid"`
	UnitPrice  float32 `json:"unitprise"`
	TotalPrice float32 `json:"totalprice"`
}
