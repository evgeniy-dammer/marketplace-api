package models

//User model
type User struct {
	Id       string `json:"id"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
