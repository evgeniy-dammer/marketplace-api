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
}
