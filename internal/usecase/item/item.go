package item

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/item"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// ItemGetAll returns all items from the system.
func (s *UseCase) ItemGetAll(ctx context.Context, userID string, organizationID string) ([]item.Item, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.ItemGetAll")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	items, err := s.adapterStorage.ItemGetAll(ctx, userID, organizationID)

	return items, errors.Wrap(err, "items select error")
}

// ItemGetOne returns item by id from the system.
func (s *UseCase) ItemGetOne(ctx context.Context, userID string, organizationID string, itemID string) (item.Item, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.ItemGetOne")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	itemSingle, err := s.adapterStorage.ItemGetOne(ctx, userID, organizationID, itemID)

	return itemSingle, errors.Wrap(err, "item select error")
}

// ItemCreate inserts item into system.
func (s *UseCase) ItemCreate(ctx context.Context, userID string, input item.CreateItemInput) (string, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.ItemCreate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	itemID, err := s.adapterStorage.ItemCreate(ctx, userID, input)

	return itemID, errors.Wrap(err, "item create error")
}

// ItemUpdate updates item by id in the system.
func (s *UseCase) ItemUpdate(ctx context.Context, userID string, input item.UpdateItemInput) error {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.ItemUpdate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.ItemUpdate(ctx, userID, input)

	return errors.Wrap(err, "item update error")
}

// ItemDelete deletes item by id from the system.
func (s *UseCase) ItemDelete(ctx context.Context, userID string, organizationID string, itemID string) error {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.ItemDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.ItemDelete(ctx, userID, organizationID, itemID)

	return errors.Wrap(err, "item delete error")
}