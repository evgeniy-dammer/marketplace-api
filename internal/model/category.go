package model

// Category model.
type Category struct {
	ID             string `json:"id" db:"id"`
	NameTm         string `json:"nametm" db:"name_tm" binding:"required"`
	NameRu         string `json:"nameru" db:"name_ru" binding:"required"`
	NameTr         string `json:"nametr" db:"name_tr" binding:"required"`
	NameEn         string `json:"nameen" db:"name_en" binding:"required"`
	Parent         string `json:"parent" db:"parent_id"`
	Level          int    `json:"level" db:"level"`
	OrganizationID string `json:"organization" db:"organization_id" binding:"required"`
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
