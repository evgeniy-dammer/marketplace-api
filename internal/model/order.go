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

// UpdateOrderInput is an input data for updating order.
type UpdateOrderInput struct {
	ID             *string      `json:"id"`
	TableID        *string      `json:"table"`
	Status         *int         `json:"status"`
	OrganizationID *string      `json:"organization"`
	TotalSum       *float32     `json:"totalsum"`
	Items          *[]OrderItem `json:"orderitems"`
}

// Validate checks if update input is nil.
func (i UpdateOrderInput) Validate() error {
	if i.ID == nil && i.TableID == nil && i.OrganizationID == nil && i.TotalSum == nil && i.Items == nil {
		return errStructHasNoValues
	}

	return nil
}
