package category

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// Category entity.
type Category struct {
	// Category ID
	ID string `json:"id" db:"id"`
	// Name Turkmen
	NameTm string `json:"nametm" db:"name_tm" binding:"required"`
	// Name Russian
	NameRu string `json:"nameru" db:"name_ru" binding:"required"`
	// Name Turkish
	NameTr string `json:"nametr" db:"name_tr" binding:"required"`
	// Name English
	NameEn string `json:"nameen" db:"name_en" binding:"required"`
	// Parent category ID
	Parent string `json:"parent" db:"parent_id"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Depth level
	Level int `json:"level" db:"level"`
}

// CreateCategoryInput entity.
type CreateCategoryInput struct {
	// Name Turkmen
	NameTm string `json:"nametm" db:"name_tm" binding:"required"`
	// Name Russian
	NameRu string `json:"nameru" db:"name_ru" binding:"required"`
	// Name Turkish
	NameTr string `json:"nametr" db:"name_tr" binding:"required"`
	// Name English
	NameEn string `json:"nameen" db:"name_en" binding:"required"`
	// Parent category ID
	Parent string `json:"parent" db:"parent_id"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Depth level
	Level int `json:"level" db:"level"`
}

// UpdateCategoryInput is an input data for updating category entity.
type UpdateCategoryInput struct {
	// Category ID
	ID *string `json:"id"`
	// Name Turkmen
	NameTm *string `json:"nametm"`
	// Name Russian
	NameRu *string `json:"nameru"`
	// Name Turkish
	NameTr *string `json:"nametr"`
	// Name English
	NameEn *string `json:"nameen"`
	// Parent category ID
	Parent *string `json:"parent"`
	// Organization ID
	OrganizationID *string `json:"organisation"`
	// Depth level
	Level *int `json:"level"`
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
		return ErrStructHasNoValues
	}

	return nil
}
