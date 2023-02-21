package model

// Item model.
type Item struct {
	ID             string  `json:"id" db:"id"`
	NameTm         string  `json:"nametm" db:"name_tm" binding:"required"`
	NameRu         string  `json:"nameru" db:"name_ru" binding:"required"`
	NameTr         string  `json:"nametr" db:"name_tr" binding:"required"`
	NameEn         string  `json:"nameen" db:"name_en" binding:"required"`
	Price          float32 `json:"price" db:"price"`
	CategoryID     string  `json:"category" db:"category_id" binding:"required"`
	OrganizationID string  `json:"organization" db:"organization_id" binding:"required"`
	Images         []Image `json:"images"`
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
