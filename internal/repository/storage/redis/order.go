package redis

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/order"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

// OrderGetAll gets orders from cache.
func (r *Repository) OrderGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]order.Order, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrderGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	orders := &order.ListOrder{}

	bytes, err := r.client.Get(ctx, ordersKey+"o."+meta.OrganizationID).Bytes()
	if err != nil {
		return *orders, errors.Wrap(err, "unable to get orders from cache")
	}

	if err = easyjson.Unmarshal(bytes, orders); err != nil {
		return *orders, errors.Wrap(err, "unable to unmarshal")
	}

	return *orders, nil
}

// OrderSetAll sets orders into cache.
func (r *Repository) OrderSetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter, orders []order.Order) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrderSetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	orderSlice := order.ListOrder(orders)

	bytes, err := easyjson.Marshal(orderSlice)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, ordersKey+"o."+meta.OrganizationID, bytes, r.options.Ttl).Err()

	return err
}

// OrderGetOne gets order by id from cache.
func (r *Repository) OrderGetOne(ctxr context.Context, orderID string) (order.Order, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrderGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr order.Order

	bytes, err := r.client.Get(ctx, orderKey+orderID).Bytes()
	if err != nil {
		return usr, errors.Wrap(err, "unable to get order from cache")
	}

	if err = easyjson.Unmarshal(bytes, &usr); err != nil {
		return usr, errors.Wrap(err, "unable to unmarshal")
	}

	return usr, nil
}

// OrderCreate sets order into cache.
func (r *Repository) OrderCreate(ctxr context.Context, usr order.Order) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrderCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, orderKey+usr.ID, bytes, r.options.Ttl).Err()

	return err
}

// OrderUpdate updates order by id in cache.
func (r *Repository) OrderUpdate(ctxr context.Context, usr order.Order) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrderUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	r.client.Set(ctx, orderKey+usr.ID, bytes, r.options.Ttl)

	return nil
}

// OrderDelete deletes order by id from cache.
func (r *Repository) OrderDelete(ctxr context.Context, orderID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrderDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, orderKey+orderID).Err()

	return err
}

// OrderInvalidate invalidate orders cache.
func (r *Repository) OrderInvalidate(ctxr context.Context) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrderInvalidate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	iter := r.client.Scan(ctx, 0, ordersKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}

	return iter.Err()
}
