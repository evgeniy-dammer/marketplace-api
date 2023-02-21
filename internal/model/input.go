package model

import "errors"

var errStructHasNoValues = errors.New("update structure has no values")

// SignInInput is an input data for signing in.
type SignInInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserInput is an input data for updating user.
type UpdateUserInput struct {
	ID        *string `json:"id"`
	FirstName *string `json:"firstname"`
	LastName  *string `json:"lastname"`
	Password  *string `json:"password"`
}

// Validate checks if update input is nil.
func (i UpdateUserInput) Validate() error {
	if i.ID == nil && i.FirstName == nil && i.LastName == nil {
		return errStructHasNoValues
	}

	return nil
}

// UpdateOrganizationInput is an input data for updating organization.
type UpdateOrganizationInput struct {
	ID      *string `json:"id"`
	Name    *string `json:"name"`
	Address *string `json:"address"`
	Phone   *string `json:"phone"`
}

// Validate checks if update input is nil.
func (i UpdateOrganizationInput) Validate() error {
	if i.ID == nil && i.Name == nil {
		return errStructHasNoValues
	}

	return nil
}

// UpdateCategoryInput is an input data for updating category.
type UpdateCategoryInput struct {
	ID             *string `json:"id"`
	NameTm         *string `json:"nametm"`
	NameRu         *string `json:"nameru"`
	NameTr         *string `json:"nametr"`
	NameEn         *string `json:"nameen"`
	Parent         *string `json:"parent"`
	Level          *int    `json:"level"`
	OrganizationID *string `json:"organisation"`
}

// Validate checks if update input is nil.
func (i UpdateCategoryInput) Validate() error {
	if i.ID == nil &&
		i.NameTm == nil &&
		i.NameRu == nil &&
		i.NameTr == nil &&
		i.NameEn == nil &&
		i.Level == nil &&
		i.OrganizationID == nil {
		return errStructHasNoValues
	}

	return nil
}

// UpdateItemInput is an input data for updating item.
type UpdateItemInput struct {
	ID             *string  `json:"id"`
	NameTm         *string  `json:"nametm"`
	NameRu         *string  `json:"nameru"`
	NameTr         *string  `json:"nametr"`
	NameEn         *string  `json:"nameen"`
	Price          *float32 `json:"price"`
	CategoryID     *string  `json:"category"`
	OrganizationID *string  `json:"organisation"`
}

// Validate checks if update input is nil.
func (i UpdateItemInput) Validate() error {
	if i.ID == nil &&
		i.NameTm == nil &&
		i.NameRu == nil &&
		i.NameTr == nil &&
		i.NameEn == nil &&
		i.CategoryID == nil &&
		i.OrganizationID == nil {
		return errStructHasNoValues
	}

	return nil
}

// UpdateTableInput is an input data for updating table.
type UpdateTableInput struct {
	ID             *string `json:"id"`
	Name           *string `json:"name"`
	OrganizationID *string `json:"organisation"`
}

// Validate checks if update input is nil.
func (i UpdateTableInput) Validate() error {
	if i.ID == nil && i.Name == nil && i.OrganizationID == nil {
		return errStructHasNoValues
	}

	return nil
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

// UpdateImageInput is an input data for updating image.
type UpdateImageInput struct {
	ID             *string `json:"id"`
	ObjectID       *string `json:"object"`
	ObjectType     *int    `json:"type"`
	Origin         *string `json:"origin"`
	Middle         *string `json:"middle"`
	Small          *string `json:"small"`
	OrganizationID *string `json:"organization"`
}

// Validate checks if update input is nil.
func (i UpdateImageInput) Validate() error {
	if i.ID == nil && i.ObjectID == nil && i.ObjectType == nil && i.OrganizationID == nil {
		return errStructHasNoValues
	}

	return nil
}
