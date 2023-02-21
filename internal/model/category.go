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
