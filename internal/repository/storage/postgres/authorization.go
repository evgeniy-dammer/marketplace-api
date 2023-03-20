package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// AuthorizationGetUserRole returns users role name.
func (r *Repository) AuthorizationGetUserRole(ctxr context.Context, userID string) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.AuthorizationGetUserRole")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var name string

	builder := r.genSQL.Select("ro.name AS role").
		From(roleTable + " ro").
		InnerJoin(userRoleTable + " ur ON ur.role_id = ro.id").
		InnerJoin(userTable + " us ON us.id = ur.user_id").
		Where(squirrel.Eq{"us.id": userID})

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.GetContext(ctx, &name, qry, args...)

	return name, errors.Wrap(err, "role name select error")
}
