package postgres

import (
	"fmt"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/favorite"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// FavoriteCreate insert favorite into database.
func (r *Repository) FavoriteCreate(ctxr context.Context, userID string, favorite favorite.Favorite) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.FavoriteCreate")
		defer span.Finish()

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

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.FavoriteDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND item_id = $2", favoriteTable)

	_, err := r.database.ExecContext(ctx, query, userID, itemID)

	return errors.Wrap(err, "favorite delete query error")
}
