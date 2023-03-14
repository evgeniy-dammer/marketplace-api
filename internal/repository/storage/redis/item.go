package redis

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/item"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

// ItemGetAll gets items from cache.
func (r *Repository) ItemGetAll(ctxr context.Context, organizationID string) ([]item.Item, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ItemGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	items := &item.ListItem{}

	bytes, err := r.client.Get(ctx, itemsKey+"o."+organizationID).Bytes()
	if err != nil {
		return *items, errors.Wrap(err, "unable to get items from cache")
	}

	if err = easyjson.Unmarshal(bytes, items); err != nil {
		return *items, errors.Wrap(err, "unable to unmarshal")
	}

	return *items, nil
}

// ItemSetAll sets items into cache.
func (r *Repository) ItemSetAll(ctxr context.Context, organizationID string, items []item.Item) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ItemSetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	itemSlice := item.ListItem(items)

	bytes, err := easyjson.Marshal(itemSlice)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, itemsKey+"o."+organizationID, bytes, r.options.Ttl).Err()

	return err
}

// ItemGetOne gets item by id from cache.
func (r *Repository) ItemGetOne(ctxr context.Context, itemID string) (item.Item, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ItemGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr item.Item

	bytes, err := r.client.Get(ctx, itemKey+itemID).Bytes()
	if err != nil {
		return usr, errors.Wrap(err, "unable to get item from cache")
	}

	if err = easyjson.Unmarshal(bytes, &usr); err != nil {
		return usr, errors.Wrap(err, "unable to unmarshal")
	}

	return usr, nil
}

// ItemCreate sets item into cache.
func (r *Repository) ItemCreate(ctxr context.Context, usr item.Item) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ItemCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, itemKey+usr.ID, bytes, r.options.Ttl).Err()

	return err
}

// ItemUpdate updates item by id in cache.
func (r *Repository) ItemUpdate(ctxr context.Context, usr item.Item) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ItemUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	r.client.Set(ctx, itemKey+usr.ID, bytes, r.options.Ttl)

	return nil
}

// ItemDelete deletes item by id from cache.
func (r *Repository) ItemDelete(ctxr context.Context, itemID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ItemDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, itemKey+itemID).Err()

	return err
}

// ItemInvalidate invalidate items cache.
func (r *Repository) ItemInvalidate(ctxr context.Context) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ItemInvalidate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	iter := r.client.Scan(ctx, 0, itemsKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}

	return iter.Err()
}
