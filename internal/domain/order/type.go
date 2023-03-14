package order

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// ListOrder
//
//easyjson:json
type ListOrder []Order

// Order entity.
//
//easyjson:json
type Order struct {
	// Order ID
	ID string `json:"id" db:"id"`
	// User ID
	UserID string `json:"user,omitempty" db:"user_id"`
	// Table ID
	TableID string `json:"table,omitempty" db:"table_id"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Created datetime
	CreatedAt string `json:"created,omitempty" db:"created_at"`
	// Updated datetime
	UpdatedAt string `json:"updated,omitempty" db:"updated_at"`
	// Order items
	Items []OrderItem `json:"orderitems,omitempty"`
	// Order status ID
	StatusID int `json:"status,omitempty" db:"status_id"`
	// Order total sum
	TotalSum float32 `json:"totalsum" db:"totalsum" binding:"required"`
}

// CreateOrderInput entity.
//
//easyjson:json
type CreateOrderInput struct {
	// User ID
	UserID string `json:"user,omitempty" db:"user_id"`
	// Table ID
	TableID string `json:"table,omitempty" db:"table_id"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Order items
	Items []CreateOrderItemInput `json:"orderitems,omitempty"`
	// Order status ID
	StatusID int `json:"status,omitempty" db:"status_id"`
	// Order total sum
	TotalSum float32 `json:"totalsum" db:"totalsum" binding:"required"`
}

// OrderItem entity.
//
//easyjson:json
type OrderItem struct {
	// OrderItem ID
	ID string `json:"id" db:"id"`
	// Order ID
	OrderID string `json:"order" db:"order_id"`
	// Item ID
	ItemID string `json:"item" db:"item_id" binding:"required"`
	// Item quantity
	Quantity float32 `json:"quantity" db:"quantity" binding:"required"`
	// Item unit price
	UnitPrice float32 `json:"unitprise" db:"unitprise" binding:"required"`
	// Item total price
	TotalPrice float32 `json:"totalprice" db:"totalprice" binding:"required"`
}

// CreateOrderItemInput entity.
//
//easyjson:json
type CreateOrderItemInput struct {
	// Item ID
	ItemID string `json:"item" db:"item_id" binding:"required"`
	// Item quantity
	Quantity float32 `json:"quantity" db:"quantity" binding:"required"`
	// Item unit price
	UnitPrice float32 `json:"unitprise" db:"unitprise" binding:"required"`
	// Item total price
	TotalPrice float32 `json:"totalprice" db:"totalprice" binding:"required"`
}

// UpdateOrderInput is an input data for updating order entity.
//
//easyjson:json
type UpdateOrderInput struct {
	// Order ID
	ID *string `json:"id"`
	// Table ID
	TableID *string `json:"table"`
	// Organization ID
	OrganizationID *string `json:"organization"`
	// Order items
	Items *[]OrderItem `json:"orderitems"`
	// Order status ID
	Status *int `json:"status"`
	// Order total sum
	TotalSum *float32 `json:"totalsum"`
}

// Validate checks if update input is nil.
func (i UpdateOrderInput) Validate() error {
	if i.ID == nil && i.TableID == nil && i.OrganizationID == nil && i.TotalSum == nil && i.Items == nil {
		return ErrStructHasNoValues
	}

	return nil
}
