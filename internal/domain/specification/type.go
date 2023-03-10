package specification

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// Specification entity.
type Specification struct {
	// Specification ID
	ID string `json:"id" db:"id"`
	// Item ID
	ItemID string `json:"item" db:"item_id" binding:"required"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Name Turkmen
	NameTm string `json:"nametm" db:"name_tm" binding:"required"`
	// Name Russian
	NameRu string `json:"nameru" db:"name_ru" binding:"required"`
	// Name Turkish
	NameTr string `json:"nametr" db:"name_tr" binding:"required"`
	// Name English
	NameEn string `json:"nameen" db:"name_en" binding:"required"`
	// Description Turkmen
	DescriptionTm string `json:"descriptiontm" db:"description_tm"`
	// Description Russian
	DescriptionRu string `json:"descriptionru" db:"description_ru"`
	// Description Turkish
	DescriptionTr string `json:"descriptiontr" db:"description_tr"`
	// Description English
	DescriptionEn string `json:"descriptionen" db:"description_en"`
	// Value
	Value string `json:"value" db:"value" binding:"required"`
}

// CreateSpecificationInput entity.
type CreateSpecificationInput struct {
	// Item ID
	ItemID string `json:"item" db:"item_id" binding:"required"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Name Turkmen
	NameTm string `json:"nametm" db:"name_tm" binding:"required"`
	// Name Russian
	NameRu string `json:"nameru" db:"name_ru" binding:"required"`
	// Name Turkish
	NameTr string `json:"nametr" db:"name_tr" binding:"required"`
	// Name English
	NameEn string `json:"nameen" db:"name_en" binding:"required"`
	// Description Turkmen
	DescriptionTm string `json:"descriptiontm" db:"description_tm"`
	// Description Russian
	DescriptionRu string `json:"descriptionru" db:"description_ru"`
	// Description Turkish
	DescriptionTr string `json:"descriptiontr" db:"description_tr"`
	// Description English
	DescriptionEn string `json:"descriptionen" db:"description_en"`
	// Value
	Value string `json:"value" db:"value" binding:"required"`
}

// UpdateSpecificationInput is an input data for updating specification entity.
type UpdateSpecificationInput struct {
	// Specification ID
	ID *string `json:"id"`
	// Item ID
	ItemID *string `json:"item"`
	// Organization ID
	OrganizationID *string `json:"organization"`
	// Name Turkmen
	NameTm *string `json:"nametm"`
	// Name Russian
	NameRu *string `json:"nameru"`
	// Name Turkish
	NameTr *string `json:"nametr"`
	// Name English
	NameEn *string `json:"nameen"`
	// Description Turkmen
	DescriptionTm *string `json:"descriptiontm"`
	// Description Russian
	DescriptionRu *string `json:"descriptionru"`
	// Description Turkish
	DescriptionTr *string `json:"descriptiontr"`
	// Description English
	DescriptionEn *string `json:"descriptionen"`
	// Value
	Value *string `json:"value"`
}

// Validate checks if update input is nil.
func (i UpdateSpecificationInput) Validate() error {
	if i.ID == nil &&
		i.ItemID == nil &&
		i.NameTm == nil &&
		i.NameRu == nil &&
		i.NameTr == nil &&
		i.NameEn == nil &&
		i.Value == nil &&
		i.OrganizationID == nil {
		return ErrStructHasNoValues
	}

	return nil
}
