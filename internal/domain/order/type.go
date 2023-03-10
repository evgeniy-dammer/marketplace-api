package order

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// Order entity.
type Order struct {
	ID             string      `json:"id" db:"id"`
	UserID         string      `json:"user,omitempty" db:"user_id"`
	TableID        string      `json:"table,omitempty" db:"table_id"`
	OrganizationID string      `json:"organization" db:"organization_id" binding:"required"`
	CreatedAt      string      `json:"created,omitempty" db:"created_at"`
	UpdatedAt      string      `json:"updated,omitempty" db:"updated_at"`
	Items          []OrderItem `json:"orderitems,omitempty"`
	StatusID       int         `json:"status,omitempty" db:"status_id"`
	TotalSum       float32     `json:"totalsum" db:"totalsum" binding:"required"`
}

// CreateOrderInput entity.
type CreateOrderInput struct {
	UserID         string                 `json:"user,omitempty" db:"user_id"`
	TableID        string                 `json:"table,omitempty" db:"table_id"`
	OrganizationID string                 `json:"organization" db:"organization_id" binding:"required"`
	Items          []CreateOrderItemInput `json:"orderitems,omitempty"`
	StatusID       int                    `json:"status,omitempty" db:"status_id"`
	TotalSum       float32                `json:"totalsum" db:"totalsum" binding:"required"`
}

// OrderItem entity.
type OrderItem struct {
	ID         string  `json:"id" db:"id"`
	OrderID    string  `json:"order" db:"order_id"`
	ItemID     string  `json:"item" db:"item_id" binding:"required"`
	Quantity   float32 `json:"quantity" db:"quantity" binding:"required"`
	UnitPrice  float32 `json:"unitprise" db:"unitprise" binding:"required"`
	TotalPrice float32 `json:"totalprice" db:"totalprice" binding:"required"`
}

// CreateOrderItemInput entity.
type CreateOrderItemInput struct {
	ItemID     string  `json:"item" db:"item_id" binding:"required"`
	Quantity   float32 `json:"quantity" db:"quantity" binding:"required"`
	UnitPrice  float32 `json:"unitprise" db:"unitprise" binding:"required"`
	TotalPrice float32 `json:"totalprice" db:"totalprice" binding:"required"`
}

// UpdateOrderInput is an input data for updating order entity.
type UpdateOrderInput struct {
	ID             *string      `json:"id"`
	TableID        *string      `json:"table"`
	OrganizationID *string      `json:"organization"`
	Items          *[]OrderItem `json:"orderitems"`
	Status         *int         `json:"status"`
	TotalSum       *float32     `json:"totalsum"`
}

// Validate checks if update input is nil.
func (i UpdateOrderInput) Validate() error {
	if i.ID == nil && i.TableID == nil && i.OrganizationID == nil && i.TotalSum == nil && i.Items == nil {
		return ErrStructHasNoValues
	}

	return nil
}
