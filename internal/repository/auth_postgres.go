package repository

import (
	"fmt"
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
)

// AuthPostgres repository
type AuthPostgres struct {
	db *sqlx.DB
}

// NewAuthPostgres constructor for AuthPostgres
func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// GetUser returns users id by username and password hash
func (r *AuthPostgres) GetUser(id string, username string) (model.User, error) {
	var user model.User
	var where string

	if id != "" {
		where = fmt.Sprint("WHERE us.id = '", id, "' ")
	} else {
		where = fmt.Sprint("WHERE us.phone = '", username, "' ")
	}

	query := fmt.Sprintf(
		"SELECT us.id, us.phone, us.password, us.first_name, us.last_name, ro.name AS role, st.name AS status FROM %s us "+
			"INNER JOIN %s ur ON ur.user_id = us.id "+
			"INNER JOIN %s ro ON ur.role_id = ro.id "+
			"INNER JOIN %s st ON st.id = us.status_id "+
			" %s ",
		userTable, userRoleTable, roleTable, statusTable, where)
	err := r.db.Get(&user, query)

	return user, err
}

// CreateUser insert User into database
func (r *AuthPostgres) CreateUser(user model.User, statusId string) (string, error) {
	var id string

	tx, err := r.db.Begin()

	if err != nil {
		return "", err
	}

	createUserQuery := fmt.Sprintf(
		"INSERT INTO %s (phone, password, first_name, last_name, status_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		userTable,
	)
	row := tx.QueryRow(createUserQuery, user.Phone, user.Password, user.FirstName, user.LastName, statusId)

	if err = row.Scan(&id); err != nil {
		if err = tx.Rollback(); err != nil {
			return "", err
		}
		return "", err
	}

	createUsersRoleQuery := fmt.Sprintf("INSERT INTO %s (user_id, role_id) VALUES ($1, $2)", userRoleTable)

	if _, err = tx.Exec(createUsersRoleQuery, id, user.RoleId); err != nil {
		if err = tx.Rollback(); err != nil {
			return "", err
		}
		return "", err
	}

	return id, tx.Commit()
}

func (r *AuthPostgres) GetUserRole(id string) (string, error) {
	var name string

	query := fmt.Sprintf("SELECT ro.name AS role FROM %s ro "+
		"INNER JOIN %s ur ON ur.role_id = ro.id "+
		"INNER JOIN %s us ON us.id = ur.user_id "+
		"WHERE us.id = '%s'",
		roleTable, userRoleTable, userTable, id)

	err := r.db.Get(&name, query)

	return name, err
}
