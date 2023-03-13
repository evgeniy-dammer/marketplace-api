package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/role"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
)

// UserGetAll selects all users from database.
func (r *Repository) UserGetAll(ctxr context.Context, search string, status string, roleID string) ([]user.User, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.UserGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var users []user.User

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
			"%s %s %s AND us.is_deleted = false",
		userTable, userRoleTable, roleTable, statusTable, search, roleID, status,
	)

	err := r.database.SelectContext(ctx, &users, query)

	return users, errors.Wrap(err, "users select query error")
}

// UserGetAllRoles selects all user roles from database.
func (r *Repository) UserGetAllRoles(ctxr context.Context) ([]role.Role, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.UserGetAllRoles")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var roles []role.Role

	query := fmt.Sprintf("SELECT id, name FROM %s ", roleTable)

	err := r.database.SelectContext(ctx, &roles, query)

	return roles, errors.Wrap(err, "roles select query error")
}

// UserGetOne select user by id from database.
func (r *Repository) UserGetOne(ctxr context.Context, userID string) (user.User, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.UserGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr user.User

	query := fmt.Sprintf(
		"SELECT us.id, us.phone, us.first_name, us.last_name, ro.name AS role, st.name AS status FROM %s us "+
			"INNER JOIN %s ur ON ur.user_id = us.id "+
			"INNER JOIN %s ro ON ur.role_id = ro.id "+
			"INNER JOIN %s st ON st.id = us.status_id "+
			"WHERE is_deleted = false AND us.id = $1", userTable, userRoleTable, roleTable, statusTable)
	err := r.database.GetContext(ctx, &usr, query, userID)

	return usr, errors.Wrap(err, "user select query error")
}

// UserCreate insert user into database.
func (r *Repository) UserCreate(ctxr context.Context, userID string, input user.CreateUserInput) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.UserCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var insertID string

	trx, err := r.database.Begin()
	if err != nil {
		return "", errors.Wrap(err, "transaction begin error")
	}

	createUserQuery := fmt.Sprintf(
		"INSERT INTO %s (phone, password, first_name, last_name, user_created) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		userTable,
	)
	row := trx.QueryRowContext(ctx, createUserQuery, input.Phone, input.Password, input.FirstName, input.LastName, userID)

	if err = row.Scan(&insertID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "user rollback error")
		}

		return "", errors.Wrap(err, "user id scan error")
	}

	createUsersRoleQuery := fmt.Sprintf("INSERT INTO %s (user_id, role_id) VALUES ($1, $2)", userRoleTable)

	if _, err = trx.ExecContext(ctx, createUsersRoleQuery, insertID, input.RoleID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "role rollback error")
		}

		return "", errors.Wrap(err, "role query execution error")
	}

	return insertID, errors.Wrap(trx.Commit(), "create transaction commit error")
}

// UserUpdate updates user by id in database.
func (r *Repository) UserUpdate(ctxr context.Context, userID string, input user.UpdateUserInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.UserUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	setValues := make([]string, 0, 4)
	args := make([]interface{}, 0, 4)
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

	setValues = append(setValues, fmt.Sprintf("user_updated=$%d", argID))
	args = append(args, userID)
	argID++

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE is_deleted = false AND id = '%s'", userTable, setQuery, *input.ID)
	_, err := r.database.ExecContext(ctx, query, args...)

	return errors.Wrap(err, "user update query error")
}

// UserDelete deletes user by id from database.
func (r *Repository) UserDelete(ctxr context.Context, userID string, dUserID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.UserDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 WHERE is_deleted = false AND id = $3",
		userTable,
	)

	_, err := r.database.ExecContext(ctx, query, time.Now().Format("2006-01-02 15:04:05"), userID, dUserID)

	return errors.Wrap(err, "user delete query error")
}
