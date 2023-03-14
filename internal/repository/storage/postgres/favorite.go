package postgres

import (
	"fmt"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/favorite"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// FavoriteCreate insert favorite into database.
func (r *Repository) FavoriteCreate(ctxr context.Context, userID string, favorite favorite.Favorite) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.FavoriteCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var favoriteID string

	query := fmt.Sprintf("INSERT INTO %s (user_id, item_id) VALUES ($1, $2) RETURNING id", favoriteTable)

	row := r.database.QueryRowContext(ctx, query, userID, favorite.ItemID)
	err := row.Scan(&favoriteID)

	return errors.Wrap(err, "favorite create query error")
}

// FavoriteDelete deletes favorite by userID and itemID from database.
func (r *Repository) FavoriteDelete(ctxr context.Context, userID string, itemID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.FavoriteDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND item_id = $2", favoriteTable)

	_, err := r.database.ExecContext(ctx, query, userID, itemID)

	return errors.Wrap(err, "favorite delete query error")
}
