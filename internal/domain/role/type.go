package role

// ListRole
//
//easyjson:json
type ListRole []Role

// Role entity.
//
//easyjson:json
type Role struct {
	// Role ID
	ID string `json:"id"`
	// Role name
	Name string `json:"name"`
}
