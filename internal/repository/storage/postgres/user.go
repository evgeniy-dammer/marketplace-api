package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/role"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/user"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// UserGetAll selects all users from database.
func (r *Repository) UserGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]user.User, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.UserGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var users []user.User

	qry, args, err := r.userGetAllQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.SelectContext(ctx, &users, qry, args...)

	return users, errors.Wrap(err, "users select query error")
}

// userGetAllQuery creates sql query.
func (r *Repository) userGetAllQuery(_ query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select(
		"us.id", "us.phone", "us.first_name", "us.last_name", "ro.name AS role", "st.name AS status",
	).From(userTable + " us").
		InnerJoin(userRoleTable + " ur ON ur.user_id = us.id").
		InnerJoin(roleTable + " ro ON ur.role_id = ro.id").
		InnerJoin(statusTable + " st ON st.id = us.status_id")

	if params.Search != "" {
		search := "%" + params.Search + "%"

		builder = builder.Where(squirrel.And{
			squirrel.Eq{"us.is_deleted": false},
			squirrel.Or{
				squirrel.Like{"us.phone": search},
				squirrel.Like{"us.first_name": search},
				squirrel.Like{"us.last_name": search},
				squirrel.Like{"ro.name": search},
				squirrel.Like{"st.name": search},
			},
		})
	} else {
		builder = builder.Where(squirrel.Eq{"us.is_deleted": false})
	}

	switch {
	case !params.StartDate.IsZero() && params.EndDate.IsZero():
		builder = builder.Where(squirrel.And{
			squirrel.GtOrEq{"us.created_at": params.StartDate.Format("2006-01-02 15:04:05")},
			squirrel.LtOrEq{"us.created_at": time.Now().Format("2006-01-02 15:04:05")},
		})
	case params.StartDate.IsZero() && !params.EndDate.IsZero():
		builder = builder.Where(squirrel.And{
			squirrel.GtOrEq{"us.created_at": time.Now().Format("2006-01-02 15:04:05")},
			squirrel.LtOrEq{"us.created_at": params.EndDate.Format("2006-01-02 15:04:05")},
		})
	case !params.StartDate.IsZero() && !params.EndDate.IsZero():
		builder = builder.Where(squirrel.And{
			squirrel.GtOrEq{"us.created_at": params.StartDate.Format("2006-01-02 15:04:05")},
			squirrel.LtOrEq{"us.created_at": params.EndDate.Format("2006-01-02 15:04:05")},
		})
	}

	if len(params.Sorts) > 0 {
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortUser)...)
	} else {
		builder = builder.OrderBy("us.created_at DESC")
	}

	if params.Pagination.Limit > 0 {
		builder = builder.Limit(params.Pagination.Limit)
	}

	if params.Pagination.Offset > 0 {
		builder = builder.Offset(params.Pagination.Offset)
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", nil, errors.Wrap(err, "unable to build a query string")
	}

	return qry, args, nil
}

// UserGetAllRoles selects all user roles from database.
func (r *Repository) UserGetAllRoles(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]role.Role, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.UserGetAllRoles")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var roles []role.Role

	qry, args, err := r.userGetAllRolesQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.SelectContext(ctx, &roles, qry, args...)

	return roles, errors.Wrap(err, "roles select query error")
}

// userGetAllRolesQuery creates sql query.
func (r *Repository) userGetAllRolesQuery(_ query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select("id", "name").From(roleTable)

	if params.Search != "" {
		builder = builder.Where(squirrel.Like{"name": "%" + params.Search + "%"})
	}

	if len(params.Sorts) > 0 {
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortRole)...)
	} else {
		builder = builder.OrderBy("id ASC")
	}

	if params.Pagination.Limit > 0 {
		builder = builder.Limit(params.Pagination.Limit)
	}

	if params.Pagination.Offset > 0 {
		builder = builder.Offset(params.Pagination.Offset)
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", nil, errors.Wrap(err, "unable to build a query string")
	}

	return qry, args, nil
}

// UserGetOne select user by id from database.
func (r *Repository) UserGetOne(ctxr context.Context, _ query.MetaData, userID string) (user.User, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.UserGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr user.User

	builder := r.genSQL.Select(
		"us.id", "us.phone", "us.first_name", "us.last_name", "ro.name AS role", "st.name AS status",
	).From(userTable + " us").
		InnerJoin(userRoleTable + " ur ON ur.user_id = us.id").
		InnerJoin(roleTable + " ro ON ur.role_id = ro.id").
		InnerJoin(statusTable + " st ON st.id = us.status_id").
		Where(squirrel.Eq{"us.is_deleted": false, "us.id": userID})

	qry, args, err := builder.ToSql()
	if err != nil {
		return usr, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.GetContext(ctx, &usr, qry, args...)

	return usr, errors.Wrap(err, "user select query error")
}

// UserCreate insert user into database.
func (r *Repository) UserCreate(ctxr context.Context, meta query.MetaData, input user.CreateUserInput) (string, error) {
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

	builder := r.genSQL.Insert(userTable).
		Columns("phone", "password", "first_name", "last_name", "user_created").
		Values(input.Phone, input.Password, input.FirstName, input.LastName, meta.UserID).
		Suffix("RETURNING \"id\"")

	createUserQuery, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := trx.QueryRowContext(ctx, createUserQuery, args...)

	if err = row.Scan(&insertID); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "user rollback error")
		}

		return "", errors.Wrap(err, "user id scan error")
	}

	builderUsersRoleQuery := r.genSQL.Insert(userRoleTable).
		Columns("user_id", "role_id").
		Values(insertID, input.RoleID)

	createUsersRoleQuery, args, err := builderUsersRoleQuery.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	if _, err = trx.ExecContext(ctx, createUsersRoleQuery, args...); err != nil {
		if err = trx.Rollback(); err != nil {
			return "", errors.Wrap(err, "role rollback error")
		}

		return "", errors.Wrap(err, "role query execution error")
	}

	return insertID, errors.Wrap(trx.Commit(), "create transaction commit error")
}

// UserUpdate updates user by id in database.
func (r *Repository) UserUpdate(ctxr context.Context, meta query.MetaData, input user.UpdateUserInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.UserUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(userTable)

	if input.FirstName != nil {
		builder = builder.Set("first_name", *input.FirstName)
	}

	if input.LastName != nil {
		builder = builder.Set("last_name", *input.LastName)
	}

	if input.Password != nil {
		builder = builder.Set("password", *input.Password)
	}

	builder = builder.Set("user_updated", meta.UserID).
		Set("updated_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_deleted": false, "id": *input.ID})

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "user update query error")
}

// UserDelete deletes user by id from database.
func (r *Repository) UserDelete(ctxr context.Context, meta query.MetaData, userID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.UserDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(userTable).
		Set("is_deleted", true).
		Set("user_deleted", meta.UserID).
		Set("deleted_at", time.Now().UTC()).
		Where(squirrel.Eq{"id": userID, "is_deleted": false})

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "user delete query error")
}
