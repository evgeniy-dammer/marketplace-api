package postgres

import (
	"fmt"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// AuthenticationGetUser returns users id by username and password hash.
func (r *Repository) AuthenticationGetUser(ctxr context.Context, userID string, username string) (user.User, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.AuthenticationGetUser")
		defer span.Finish()

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

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.AuthenticationCreateUser")
		defer span.Finish()

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

// AuthenticationGetUserRole returns users role name
func (r *Repository) AuthenticationGetUserRole(ctxr context.Context, userID string) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.AuthenticationGetUserRole")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var name string

	query := fmt.Sprintf("SELECT ro.name AS role FROM %s ro "+
		"INNER JOIN %s ur ON ur.role_id = ro.id "+
		"INNER JOIN %s us ON us.id = ur.user_id "+
		"WHERE us.id = '%s'",
		roleTable, userRoleTable, userTable, userID,
	)

	err := r.database.GetContext(ctx, &name, query)

	return name, errors.Wrap(err, "role name select error")
}
