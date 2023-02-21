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

// SignInInput is an input data for signing in.
type SignInInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserInput is an input data for updating user.
type UpdateUserInput struct {
	ID        *string `json:"id"`
	FirstName *string `json:"firstname"`
	LastName  *string `json:"lastname"`
	Password  *string `json:"password"`
}

// Validate checks if update input is nil.
func (i UpdateUserInput) Validate() error {
	if i.ID == nil && i.FirstName == nil && i.LastName == nil {
		return errStructHasNoValues
	}

	return nil
}
