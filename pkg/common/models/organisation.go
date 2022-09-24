package models

type Organisation struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	UserId  string `json:"userid"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}
