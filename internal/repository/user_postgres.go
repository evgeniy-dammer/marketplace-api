package repository

import (
	"fmt"
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

// UserPostgresql repository
type UserPostgresql struct {
	db *sqlx.DB
}

// NewUserPostgresql is a constructor for UserPostgresql
func NewUserPostgresql(db *sqlx.DB) *UserPostgresql {
	return &UserPostgresql{db: db}
}

// GetAll selects all users from database
func (r *UserPostgresql) GetAll(search string, status string, roleId string) ([]model.User, error) {
	var users []model.User

	if search != "" {
		search = fmt.Sprint("WHERE (us.phone LIKE '%", search, "%' OR us.first_name LIKE '%", search, "%' OR us.first_name LIKE '%", search, "%') ")
	}

	if roleId != "" {
		if search == "" {
			roleId = fmt.Sprint("WHERE ur.role_id = ", roleId, " ")
		} else {
			roleId = fmt.Sprint(" AND ur.role_id = ", roleId, " ")
		}
	}

	if status != "" {
		if search == "" && roleId == "" {
			status = fmt.Sprint("WHERE us.status_id = '", status, "' ")
		} else {
			status = fmt.Sprint(" AND us.status_id = '", status, "' ")
		}
	}

	query := fmt.Sprintf("SELECT us.id, us.phone, us.first_name, us.last_name, ro.name AS role, st.name AS status FROM %s us "+
		"INNER JOIN %s ur ON ur.user_id = us.id "+
		"INNER JOIN %s ro ON ur.role_id = ro.id "+
		"INNER JOIN %s st ON st.id = us.status_id "+
		"%s %s %s ",
		userTable, userRoleTable, roleTable, statusTable, search, roleId, status)

	err := r.db.Select(&users, query)

	return users, err
}

// GetAllRoles selects all user roles from database
func (r *UserPostgresql) GetAllRoles() ([]model.Role, error) {
	var roles []model.Role

	query := fmt.Sprintf("SELECT id, name FROM %s ", roleTable)

	err := r.db.Select(&roles, query)

	return roles, err
}

// GetOne select user by id from database
func (r *UserPostgresql) GetOne(userId string) (model.User, error) {
	var user model.User

	query := fmt.Sprintf(
		"SELECT us.id, us.phone, us.first_name, us.last_name, ro.name AS role, st.name AS status FROM %s us "+
			"INNER JOIN %s ur ON ur.user_id = us.id "+
			"INNER JOIN %s ro ON ur.role_id = ro.id "+
			"INNER JOIN %s st ON st.id = us.status_id "+
			"WHERE us.id = $1", userTable, userRoleTable, roleTable, statusTable)
	err := r.db.Get(&user, query, userId)

	return user, err
}

// Create insert User into database
func (r *UserPostgresql) Create(user model.User, statusId string) (string, error) {
	var id string

	tx, err := r.db.Begin()

	if err != nil {
		return "", err
	}

	var createUserQuery = fmt.Sprintf(
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

// Update updates list by id in database
func (r *UserPostgresql) Update(userId string, input model.UpdateUserInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.FirstName != nil {
		setValues = append(setValues, fmt.Sprintf("first_name=$%d", argId))
		args = append(args, *input.FirstName)
		argId++
	}

	if input.LastName != nil {
		setValues = append(setValues, fmt.Sprintf("last_name=$%d", argId))
		args = append(args, *input.LastName)
		argId++
	}

	if input.Password != nil {
		setValues = append(setValues, fmt.Sprintf("password=$%d", argId))
		args = append(args, *input.Password)
		argId++
	}

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))
	argId++

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = '%s'", userTable, setQuery, userId)
	_, err := r.db.Exec(query, args...)

	return err
}

// Delete deletes user by id from database
func (r *UserPostgresql) Delete(userId string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", userTable)
	_, err := r.db.Exec(query, userId)

	return err
}

// GetActiveStatusId select status ID by name from database
func (r *UserPostgresql) GetActiveStatusId(name string) (string, error) {
	var statusId string

	query := fmt.Sprintf("SELECT id FROM %s WHERE name = '%s'", statusTable, name)

	err := r.db.Get(&statusId, query)

	return statusId, err
}
