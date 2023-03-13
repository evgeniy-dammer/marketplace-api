package postgres

import (
	"fmt"

	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
)

// AuthorizationGetUserRole returns users role name
func (r *Repository) AuthorizationGetUserRole(ctxr context.Context, userID string) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.AuthorizationGetUserRole")
		defer span.End()

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
