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
	DescriptionEn  string                        `json:"descriptionen" db:"description_en"`
	DescriptionRu  string                        `json:"descriptionru" db:"description_ru"`
	NameRu         string                        `json:"nameru" db:"name_ru" binding:"required"`
	NameTr         string                        `json:"nametr" db:"name_tr" binding:"required"`
	CreatedAt      string                        `json:"created,omitempty" db:"created_at"`
	DescriptionTm  string                        `json:"descriptiontm" db:"description_tm"`
	NameTm         string                        `json:"nametm" db:"name_tm" binding:"required"`
	DescriptionTr  string                        `json:"descriptiontr" db:"description_tr"`
	NameEn         string                        `json:"nameen" db:"name_en" binding:"required"`
	InternalID     string                        `json:"internal" db:"internal_id"`
	ID             string                        `json:"id" db:"id"`
	OrganizationID string                        `json:"organization" db:"organization_id" binding:"required"`
	CategoryID     string                        `json:"category" db:"category_id" binding:"required"`
	Images         []image.Image                 `json:"images"`
	Comments       []comment.Comment             `json:"comments"`
	Specification  []specification.Specification `json:"specification"`
	BrandID        int                           `json:"brand" db:"brand_id" binding:"required"`
	CommentsQty    int                           `json:"commentsqty" db:"comments_qty"`
	Rating         float32                       `json:"rating" db:"rating"`
	Price          float32                       `json:"price" db:"price"`
}

// CreateItemInput entity.
type CreateItemInput struct {
	DescriptionEn  string  `json:"descriptionen" db:"description_en"`
	DescriptionRu  string  `json:"descriptionru" db:"description_ru"`
	NameRu         string  `json:"nameru" db:"name_ru" binding:"required"`
	NameTr         string  `json:"nametr" db:"name_tr" binding:"required"`
	DescriptionTm  string  `json:"descriptiontm" db:"description_tm"`
	NameTm         string  `json:"nametm" db:"name_tm" binding:"required"`
	DescriptionTr  string  `json:"descriptiontr" db:"description_tr"`
	NameEn         string  `json:"nameen" db:"name_en" binding:"required"`
	InternalID     string  `json:"internal" db:"internal_id"`
	OrganizationID string  `json:"organization" db:"organization_id" binding:"required"`
	CategoryID     string  `json:"category" db:"category_id" binding:"required"`
	BrandID        int     `json:"brand" db:"brand_id" binding:"required"`
	Price          float32 `json:"price" db:"price"`
}

// UpdateItemInput is an input data for updating item entity.
type UpdateItemInput struct {
	ID             *string  `json:"id"`
	NameTm         *string  `json:"nametm"`
	NameRu         *string  `json:"nameru"`
	NameTr         *string  `json:"nametr"`
	NameEn         *string  `json:"nameen"`
	DescriptionTm  *string  `json:"descriptiontm"`
	DescriptionRu  *string  `json:"descriptionru"`
	DescriptionTr  *string  `json:"descriptiontr"`
	DescriptionEn  *string  `json:"descriptionen"`
	InternalID     *string  `json:"internal"`
	CategoryID     *string  `json:"category"`
	OrganizationID *string  `json:"organisation"`
	BrandID        *int     `json:"brand"`
	Price          *float32 `json:"price"`
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
