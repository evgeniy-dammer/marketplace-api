package model

// User model.
type User struct {
	ID        string `json:"id" db:"id"`
	Phone     string `json:"phone" db:"phone" binding:"required"`
	Password  string `json:"password,omitempty" db:"password" binding:"required"`
	FirstName string `json:"firstname" db:"first_name" binding:"required"`
	LastName  string `json:"lastname" db:"last_name" binding:"required"`
	RoleName  string `json:"role,omitempty" db:"role"`
	RoleID    int    `json:"roleid,omitempty"`
	Status    string `json:"status" db:"status"`
}
