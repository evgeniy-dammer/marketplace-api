package model

// Item model.
type Item struct {
	ID             string          `json:"id" db:"id"`
	NameTm         string          `json:"nametm" db:"name_tm" binding:"required"`
	NameRu         string          `json:"nameru" db:"name_ru" binding:"required"`
	NameTr         string          `json:"nametr" db:"name_tr" binding:"required"`
	NameEn         string          `json:"nameen" db:"name_en" binding:"required"`
	DescriptionTm  string          `json:"descriptiontm" db:"description_tm"`
	DescriptionRu  string          `json:"descriptionru" db:"description_ru"`
	DescriptionTr  string          `json:"descriptiontr" db:"description_tr"`
	DescriptionEn  string          `json:"descriptionen" db:"description_en"`
	InternalID     string          `json:"internal" db:"internal_id"`
	Price          float32         `json:"price" db:"price"`
	Rating         float32         `json:"rating" db:"rating"`
	CommentsQty    int             `json:"commentsqty" db:"comments_qty"`
	CategoryID     string          `json:"category" db:"category_id" binding:"required"`
	OrganizationID string          `json:"organization" db:"organization_id" binding:"required"`
	BrandID        int             `json:"brand" db:"brand_id" binding:"required"`
	CreatedAt      string          `json:"created,omitempty" db:"created_at"`
	Images         []Image         `json:"images"`
	Comments       []Comment       `json:"comments"`
	Specification  []Specification `json:"specification"`
}

// UpdateItemInput is an input data for updating item.
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
	Price          *float32 `json:"price"`
	CategoryID     *string  `json:"category"`
	OrganizationID *string  `json:"organisation"`
	BrandID        *int     `json:"brand"`
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
