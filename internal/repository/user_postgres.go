package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// UserPostgresql repository.
type UserPostgresql struct {
	db *sqlx.DB
}

// NewUserPostgresql is a constructor for UserPostgresql.
func NewUserPostgresql(db *sqlx.DB) *UserPostgresql {
	return &UserPostgresql{db: db}
}

// GetAll selects all users from database.
func (r *UserPostgresql) GetAll(search string, status string, roleID string) ([]model.User, error) {
	var users []model.User

	if search != "" {
		search = fmt.Sprint(
			"WHERE (us.phone LIKE '%",
			search,
			"%' OR us.first_name LIKE '%",
			search,
			"%' OR us.first_name LIKE '%",
			search,
			"%') ",
		)
	}

	if roleID != "" {
		if search == "" {
			roleID = fmt.Sprint("WHERE ur.role_id = ", roleID, " ")
		} else {
			roleID = fmt.Sprint(" AND ur.role_id = ", roleID, " ")
		}
	}

	if status != "" {
		if search == "" && roleID == "" {
			status = fmt.Sprint("WHERE us.status_id = '", status, "' ")
		} else {
			status = fmt.Sprint(" AND us.status_id = '", status, "' ")
		}
	}

	query := fmt.Sprintf(
		"SELECT us.id, us.phone, us.first_name, us.last_name, ro.name AS role, st.name AS status FROM %s us "+
			"INNER JOIN %s ur ON ur.user_id = us.id "+
			"INNER JOIN %s ro ON ur.role_id = ro.id "+
			"INNER JOIN %s st ON st.id = us.status_id "+
			"%s %s %s ",
		userTable, userRoleTable, roleTable, statusTable, search, roleID, status,
	)

	err := r.db.Select(&users, query)

	return users, errors.Wrap(err, "users select query error")
}

// GetAllRoles selects all user roles from database.
func (r *UserPostgresql) GetAllRoles() ([]model.Role, error) {
	var roles []model.Role

	query := fmt.Sprintf("SELECT id, name FROM %s ", roleTable)

	err := r.db.Select(&roles, query)

	return roles, errors.Wrap(err, "roles select query error")
}

// GetOne select user by id from database.
func (r *UserPostgresql) GetOne(userID string) (model.User, error) {
	var user model.User

	query := fmt.Sprintf(
		"SELECT us.id, us.phone, us.first_name, us.last_name, ro.name AS role, st.name AS status FROM %s us "+
			"INNER JOIN %s ur ON ur.user_id = us.id "+
			"INNER JOIN %s ro ON ur.role_id = ro.id "+
			"INNER JOIN %s st ON st.id = us.status_id "+
			"WHERE us.id = $1", userTable, userRoleTable, roleTable, statusTable)
	err := r.db.Get(&user, query, userID)

	return user, errors.Wrap(err, "user select query error")
}

// Create insert user into database.
func (r *UserPostgresql) Create(user model.User, statusID string) (string, error) {
	var userID string

	trx, err := r.db.Begin()
	if err != nil {
		return "", errors.Wrap(err, "transaction begin error")
	}

	createUserQuery := fmt.Sprintf(
		"INSERT INTO %s (phone, password, first_name, last_name, status_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		userTable,
	)
	row := trx.QueryRow(createUserQuery, user.Phone, user.Password, user.FirstName, user.LastName, statusID)

	if err = row.Scan(&userID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "user rollback error")
		}

		return "", errors.Wrap(err, "user id scan error")
	}

	createUsersRoleQuery := fmt.Sprintf("INSERT INTO %s (user_id, role_id) VALUES ($1, $2)", userRoleTable)

	if _, err = trx.Exec(createUsersRoleQuery, userID, user.RoleID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "role rollback error")
		}

		return "", errors.Wrap(err, "role query execution error")
	}

	return userID, errors.Wrap(trx.Commit(), "create transaction commit error")
}

// Update updates user by id in database.
func (r *UserPostgresql) Update(userID string, input model.UpdateUserInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.FirstName != nil {
		setValues = append(setValues, fmt.Sprintf("first_name=$%d", argID))
		args = append(args, *input.FirstName)
		argID++
	}

	if input.LastName != nil {
		setValues = append(setValues, fmt.Sprintf("last_name=$%d", argID))
		args = append(args, *input.LastName)
		argID++
	}

	if input.Password != nil {
		setValues = append(setValues, fmt.Sprintf("password=$%d", argID))
		args = append(args, *input.Password)
		argID++
	}

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = '%s'", userTable, setQuery, userID)
	_, err := r.db.Exec(query, args...)

	return errors.Wrap(err, "user update query error")
}

// Delete deletes user by id from database.
func (r *UserPostgresql) Delete(userID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", userTable)
	_, err := r.db.Exec(query, userID)

	return errors.Wrap(err, "user delete query error")
}

// GetActiveStatusID select status ID by name from database.
func (r *UserPostgresql) GetActiveStatusID(name string) (string, error) {
	var statusID string

	query := fmt.Sprintf("SELECT id FROM %s WHERE name = '%s'", statusTable, name)

	err := r.db.Get(&statusID, query)

	return statusID, errors.Wrap(err, "active status select query error")
}
