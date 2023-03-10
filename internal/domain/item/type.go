package item

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/comment"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/image"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/specification"
	"github.com/pkg/errors"
)

var ErrStructHasNoValues = errors.New("update structure has no values")

// Item entity.
type Item struct {
	// Description English
	DescriptionEn string `json:"descriptionen" db:"description_en"`
	// Description Russian
	DescriptionRu string `json:"descriptionru" db:"description_ru"`
	// Name Russian
	NameRu string `json:"nameru" db:"name_ru" binding:"required"`
	// Name Turkish
	NameTr string `json:"nametr" db:"name_tr" binding:"required"`
	// Created datetime
	CreatedAt string `json:"created,omitempty" db:"created_at"`
	// Description Turkmen
	DescriptionTm string `json:"descriptiontm" db:"description_tm"`
	// Name Turkmen
	NameTm string `json:"nametm" db:"name_tm" binding:"required"`
	// Description Turkish
	DescriptionTr string `json:"descriptiontr" db:"description_tr"`
	// Name English
	NameEn string `json:"nameen" db:"name_en" binding:"required"`
	// Internal ID
	InternalID string `json:"internal" db:"internal_id"`
	// Item ID
	ID string `json:"id" db:"id"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Category ID
	CategoryID string `json:"category" db:"category_id" binding:"required"`
	// Images
	Images []image.Image `json:"images"`
	// Comments
	Comments []comment.Comment `json:"comments"`
	// Specifications
	Specification []specification.Specification `json:"specification"`
	// Brand ID
	BrandID int `json:"brand" db:"brand_id" binding:"required"`
	// Comments quantity
	CommentsQty int `json:"commentsqty" db:"comments_qty"`
	// Rating value
	Rating float32 `json:"rating" db:"rating"`
	// Item price
	Price float32 `json:"price" db:"price"`
}

// CreateItemInput entity.
type CreateItemInput struct {
	// Description English
	DescriptionEn string `json:"descriptionen" db:"description_en"`
	// Description Russian
	DescriptionRu string `json:"descriptionru" db:"description_ru"`
	// Name Russian
	NameRu string `json:"nameru" db:"name_ru" binding:"required"`
	// Name Turkish
	NameTr string `json:"nametr" db:"name_tr" binding:"required"`
	// Description Turkmen
	DescriptionTm string `json:"descriptiontm" db:"description_tm"`
	// Name Turkmen
	NameTm string `json:"nametm" db:"name_tm" binding:"required"`
	// Description Turkish
	DescriptionTr string `json:"descriptiontr" db:"description_tr"`
	// Name English
	NameEn string `json:"nameen" db:"name_en" binding:"required"`
	// Internal ID
	InternalID string `json:"internal" db:"internal_id"`
	// Organization ID
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
	// Category ID
	CategoryID string `json:"category" db:"category_id" binding:"required"`
	// Brand ID
	BrandID int `json:"brand" db:"brand_id" binding:"required"`
	// Item price
	Price float32 `json:"price" db:"price"`
}

// UpdateItemInput is an input data for updating item entity.
type UpdateItemInput struct {
	// Item ID
	ID *string `json:"id"`
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
	// Internal ID
	InternalID *string `json:"internal"`
	// Category ID
	CategoryID *string `json:"category"`
	// Organization ID
	OrganizationID *string `json:"organisation"`
	// Brand ID
	BrandID *int `json:"brand"`
	// Item price
	Price *float32 `json:"price"`
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
		return ErrStructHasNoValues
	}

	return nil
}
