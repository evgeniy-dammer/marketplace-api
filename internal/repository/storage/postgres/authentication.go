package postgres

import (
	"fmt"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
)

// AuthenticationGetUser returns users id by username and password hash.
func (r *Repository) AuthenticationGetUser(ctxr context.Context, userID string, username string) (user.User, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.AuthenticationGetUser")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr user.User

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

	err := r.database.GetContext(ctx, &usr, query)

	return usr, errors.Wrap(err, "user select error")
}

// AuthenticationCreateUser insert user into database.
func (r *Repository) AuthenticationCreateUser(ctxr context.Context, input user.CreateUserInput) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.AuthenticationCreateUser")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var userID string

	trx, err := r.database.Begin()
	if err != nil {
		return "", errors.Wrap(err, "transaction begin error")
	}

	createUserQuery := fmt.Sprintf(
		"INSERT INTO %s (phone, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id",
		userTable,
	)

	row := trx.QueryRowContext(ctx, createUserQuery, input.Phone, input.Password, input.FirstName, input.LastName)

	if err = row.Scan(&userID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "transaction rollback error")
		}

		return "", errors.Wrap(err, "user id scan error")
	}

	createUsersRoleQuery := fmt.Sprintf("INSERT INTO %s (user_id, role_id) VALUES ($1, $2)", userRoleTable)

	if _, err = trx.ExecContext(ctx, createUsersRoleQuery, userID, input.RoleID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "role table rollback error")
		}

		return "", errors.Wrap(err, "role insert query error")
	}

	return userID, errors.Wrap(trx.Commit(), "transaction commit error")
}
