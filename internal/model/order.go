package model

// Order model.
type Order struct {
	ID             string      `json:"id" db:"id"`
	UserID         string      `json:"user,omitempty" db:"user_id"`
	TableID        string      `json:"table,omitempty" db:"table_id"`
	StatusID       int         `json:"status,omitempty" db:"status_id"`
	OrganizationID string      `json:"organization" db:"organization_id" binding:"required"`
	TotalSum       float32     `json:"totalsum" db:"totalsum" binding:"required"`
	CreatedAt      string      `json:"created,omitempty" db:"created_at"`
	UpdatedAt      string      `json:"updated,omitempty" db:"updated_at"`
	Items          []OrderItem `json:"orderitems,omitempty"`
}
