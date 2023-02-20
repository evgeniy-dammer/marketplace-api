package model

// OrderItem model.
type OrderItem struct {
	ID         string  `json:"id" db:"id"`
	OrderID    string  `json:"order" db:"order_id"`
	ItemID     string  `json:"item" db:"item_id" binding:"required"`
	Quantity   float32 `json:"quantity" db:"quantity" binding:"required"`
	UnitPrice  float32 `json:"unitprise" db:"unitprise" binding:"required"`
	TotalPrice float32 `json:"totalprice" db:"totalprice" binding:"required"`
}
