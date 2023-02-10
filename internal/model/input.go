package model

import "errors"

// SignInInput is an input data for signing in
type SignInInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserInput is an input data for updating user
type UpdateUserInput struct {
	FirstName *string `json:"firstname"`
	LastName  *string `json:"lastname"`
	Password  *string `json:"password"`
}

// Validate checks if update input is nil
func (i UpdateUserInput) Validate() error {
	if i.FirstName == nil && i.LastName == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

// UpdateOrganizationInput is an input data for updating organization
type UpdateOrganizationInput struct {
	Name    *string `json:"name"`
	Address *string `json:"address"`
	Phone   *string `json:"phone"`
}

// Validate checks if update input is nil
func (i UpdateOrganizationInput) Validate() error {
	if i.Name == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

// UpdateCategoryInput is an input data for updating category
type UpdateCategoryInput struct {
	Name   *string `json:"name"`
	Parent *string `json:"parent"`
	Level  *int    `json:"level"`
}

// Validate checks if update input is nil
func (i UpdateCategoryInput) Validate() error {
	if i.Name == nil && i.Level == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

// UpdateItemInput is an input data for updating organization
type UpdateItemInput struct {
	Name           *string  `json:"name"`
	Price          *float32 `json:"price"`
	CategoryId     *string  `json:"category"`
	OrganisationId *string  `json:"organisation"`
}

// Validate checks if update input is nil
func (i UpdateItemInput) Validate() error {
	if i.Name == nil && i.OrganisationId == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
