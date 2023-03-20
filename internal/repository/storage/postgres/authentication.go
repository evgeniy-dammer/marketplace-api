package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/user"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
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

	builder := r.genSQL.Select(
		"us.id", "us.phone", "us.password", "us.first_name", "us.last_name", "ro.name AS role", "st.name AS status").
		From(userTable + " us").
		InnerJoin(userRoleTable + " ur ON ur.user_id = us.id").
		InnerJoin(roleTable + " ro ON ur.role_id = ro.id").
		InnerJoin(statusTable + " st ON st.id = us.status_id").
		Where(squirrel.Eq{"us.is_deleted": false})

	if userID != "" {
		builder = builder.Where(squirrel.Eq{"us.id": userID})
	} else {
		builder = builder.Where(squirrel.Eq{"us.phone": username})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return usr, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.GetContext(ctx, &usr, qry, args...)

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

	builder := r.genSQL.Insert(userTable).
		Columns("phone", "password", "first_name", "last_name").
		Values(input.Phone, input.Password, input.FirstName, input.LastName).
		Suffix("RETURNING \"id\"")

	createUserQuery, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := trx.QueryRowContext(ctx, createUserQuery, args...)

	if err = row.Scan(&userID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "transaction rollback error")
		}

		return "", errors.Wrap(err, "user id scan error")
	}

	builderUsersRoleQuery := r.genSQL.Insert(userRoleTable).
		Columns("user_id", "role_id").
		Values(userID, input.RoleID)

	createUsersRoleQuery, args, err := builderUsersRoleQuery.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	if _, err = trx.ExecContext(ctx, createUsersRoleQuery, args...); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "role table rollback error")
		}

		return "", errors.Wrap(err, "role insert query error")
	}

	return userID, errors.Wrap(trx.Commit(), "transaction commit error")
}

// AuthenticationCreateToken inserts token into database.
func (r *Repository) AuthenticationCreateToken(ctxr context.Context, userID string, token string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.AuthenticationCreateToken")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var tokenID string

	builder := r.genSQL.Insert(tokenTable).Columns("user_id", "token").Values(userID, token).Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	row := r.database.QueryRowContext(ctx, qry, args...)

	err = row.Scan(&tokenID)
	if err != nil {
		return errors.Wrap(err, "unable to scan token id")
	}

	return nil
}

// AuthenticationGetToken returns token from database.
func (r *Repository) AuthenticationGetToken(ctxr context.Context, userID string, token string) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.AuthenticationGetToken")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var tokenID string

	builder := r.genSQL.Select("id").From(tokenTable).
		Where(squirrel.Eq{"user_id": userID, "token": token, "expired": false})

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.GetContext(ctx, &tokenID, qry, args...)
	if err != nil {
		return "", errors.Wrap(err, "unable to create select token id")
	}

	return tokenID, nil
}

// AuthenticationUpdateToken updates token in database
func (r *Repository) AuthenticationUpdateToken(ctxr context.Context, tokenID string, token string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.AuthenticationGetToken")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(tokenTable).
		Set("token", token).
		Set("expired", false).
		Where(squirrel.Eq{"id": tokenID})

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)
	if err != nil {
		return errors.Wrap(err, "unable to update token")
	}

	return nil
}
