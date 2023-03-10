package specification

import (
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// Specification entity.
type Specification struct {
	ID             string `json:"id" db:"id"`
	ItemID         string `json:"item" db:"item_id" binding:"required"`
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	NameTm         string `json:"nametm" db:"name_tm" binding:"required"`
	NameRu         string `json:"nameru" db:"name_ru" binding:"required"`
	NameTr         string `json:"nametr" db:"name_tr" binding:"required"`
	NameEn         string `json:"nameen" db:"name_en" binding:"required"`
	DescriptionTm  string `json:"descriptiontm" db:"description_tm"`
	DescriptionRu  string `json:"descriptionru" db:"description_ru"`
	DescriptionTr  string `json:"descriptiontr" db:"description_tr"`
	DescriptionEn  string `json:"descriptionen" db:"description_en"`
	Value          string `json:"value" db:"value" binding:"required"`
}

// CreateSpecificationInput entity.
type CreateSpecificationInput struct {
	ItemID         string `json:"item" db:"item_id" binding:"required"`
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	NameTm         string `json:"nametm" db:"name_tm" binding:"required"`
	NameRu         string `json:"nameru" db:"name_ru" binding:"required"`
	NameTr         string `json:"nametr" db:"name_tr" binding:"required"`
	NameEn         string `json:"nameen" db:"name_en" binding:"required"`
	DescriptionTm  string `json:"descriptiontm" db:"description_tm"`
	DescriptionRu  string `json:"descriptionru" db:"description_ru"`
	DescriptionTr  string `json:"descriptiontr" db:"description_tr"`
	DescriptionEn  string `json:"descriptionen" db:"description_en"`
	Value          string `json:"value" db:"value" binding:"required"`
}

// UpdateSpecificationInput is an input data for updating specification entity.
type UpdateSpecificationInput struct {
	ID             *string `json:"id"`
	ItemID         *string `json:"item"`
	OrganizationID *string `json:"organization"`
	NameTm         *string `json:"nametm"`
	NameRu         *string `json:"nameru"`
	NameTr         *string `json:"nametr"`
	NameEn         *string `json:"nameen"`
	DescriptionTm  *string `json:"descriptiontm"`
	DescriptionRu  *string `json:"descriptionru"`
	DescriptionTr  *string `json:"descriptiontr"`
	DescriptionEn  *string `json:"descriptionen"`
	Value          *string `json:"value"`
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
