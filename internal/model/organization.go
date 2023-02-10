package model

// Organization model
type Organization struct {
	Id      string `json:"id" db:"id"`
	Name    string `json:"name" db:"name" binding:"required"`
	UserId  string `json:"userid" db:"user_id" binding:"required"`
	Address string `json:"address" db:"address" binding:"required"`
	Phone   string `json:"phone" db:"phone" binding:"required"`
}
