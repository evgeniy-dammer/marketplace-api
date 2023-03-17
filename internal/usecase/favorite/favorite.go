package favorite

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/favorite"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// FavoriteCreate inserts favorite into system.
func (s *UseCase) FavoriteCreate(ctx context.Context, meta query.MetaData, favorite favorite.Favorite) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.FavoriteCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.FavoriteCreate(ctx, meta, favorite)

	return errors.Wrap(err, "favorite create error")
}

// FavoriteDelete deletes favorite by id from the system.
func (s *UseCase) FavoriteDelete(ctx context.Context, meta query.MetaData, itemID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.FavoriteDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.FavoriteDelete(ctx, meta, itemID)

	return errors.Wrap(err, "favorite delete error")
}
