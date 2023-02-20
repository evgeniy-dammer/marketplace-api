package repository

import (
	"fmt"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AuthPostgres repository.
type AuthPostgres struct {
	db *sqlx.DB
}

// NewAuthPostgres constructor for AuthPostgres.
func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// GetUser returns users id by username and password hash.
func (r *AuthPostgres) GetUser(userID string, username string) (model.User, error) {
	var user model.User

	var where string

	if userID != "" {
		where = fmt.Sprint("WHERE us.id = '", userID, "' ")
	} else {
		where = fmt.Sprint("WHERE us.phone = '", username, "' ")
	}

	query := fmt.Sprintf(
		"SELECT us.id, us.phone, us.password, us.first_name, us.last_name, ro.name AS role, st.name AS status FROM %s us "+
			"INNER JOIN %s ur ON ur.user_id = us.id "+
			"INNER JOIN %s ro ON ur.role_id = ro.id "+
			"INNER JOIN %s st ON st.id = us.status_id "+
			" %s AND us.is_deleted = false",
		userTable, userRoleTable, roleTable, statusTable, where,
	)

	err := r.db.Get(&user, query)

	return user, errors.Wrap(err, "user select error")
}

// CreateUser insert user into database.
func (r *AuthPostgres) CreateUser(user model.User) (string, error) {
	var userID string

	trx, err := r.db.Begin()
	if err != nil {
		return "", errors.Wrap(err, "transaction begin error")
	}

	createUserQuery := fmt.Sprintf(
		"INSERT INTO %s (phone, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id",
		userTable,
	)

	row := trx.QueryRow(createUserQuery, user.Phone, user.Password, user.FirstName, user.LastName)

	if err = row.Scan(&userID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "transaction rollback error")
		}

		return "", errors.Wrap(err, "user id scan error")
	}

	createUsersRoleQuery := fmt.Sprintf("INSERT INTO %s (user_id, role_id) VALUES ($1, $2)", userRoleTable)

	if _, err = trx.Exec(createUsersRoleQuery, userID, user.RoleID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "role table rollback error")
		}

		return "", errors.Wrap(err, "role insert query error")
	}

	return userID, errors.Wrap(trx.Commit(), "transaction commit error")
}

// GetUserRole returns users role name
func (r *AuthPostgres) GetUserRole(userID string) (string, error) {
	var name string

	query := fmt.Sprintf("SELECT ro.name AS role FROM %s ro "+
		"INNER JOIN %s ur ON ur.role_id = ro.id "+
		"INNER JOIN %s us ON us.id = ur.user_id "+
		"WHERE us.id = '%s'",
		roleTable, userRoleTable, userTable, userID,
	)

	err := r.db.Get(&name, query)

	return name, errors.Wrap(err, "role name select error")
}
