package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/favorite"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// FavoriteCreate insert favorite into database.
func (r *Repository) FavoriteCreate(ctxr context.Context, meta query.MetaData, favorite favorite.Favorite) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.FavoriteCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var favoriteID string

	builder := r.genSQL.Insert(favoriteTable).
		Columns("user_id", "item_id").
		Values(meta.UserID, favorite.ItemID).
		Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	row := r.database.QueryRowContext(ctx, qry, args...)
	err = row.Scan(&favoriteID)

	return errors.Wrap(err, "favorite create query error")
}

// FavoriteDelete deletes favorite by userID and itemID from database.
func (r *Repository) FavoriteDelete(ctxr context.Context, meta query.MetaData, itemID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.FavoriteDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Delete(favoriteTable).Where(squirrel.Eq{"user_id": meta.UserID, "item_id": itemID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "favorite delete query error")
}
