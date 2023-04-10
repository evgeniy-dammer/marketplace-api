package redis

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/specification"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

// SpecificationGetAll gets specifications from cache.
func (r *Repository) SpecificationGetAll(ctxr context.Context, meta query.MetaData, _ queryparameter.QueryParameter) ([]specification.Specification, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.SpecificationGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	specifications := &specification.ListSpecification{}

	bytes, err := r.client.Get(ctx, specificationsKey+"o."+meta.OrganizationID).Bytes()
	if err != nil {
		return *specifications, errors.Wrap(err, "unable to get specifications from cache")
	}

	if err = easyjson.Unmarshal(bytes, specifications); err != nil {
		return *specifications, errors.Wrap(err, "unable to unmarshal")
	}

	return *specifications, nil
}

// SpecificationSetAll sets specifications into cache.
func (r *Repository) SpecificationSetAll(ctxr context.Context, meta query.MetaData, _ queryparameter.QueryParameter, specifications []specification.Specification) error { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.SpecificationSetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	specificationSlice := specification.ListSpecification(specifications)

	bytes, err := easyjson.Marshal(specificationSlice)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, specificationsKey+"o."+meta.OrganizationID, bytes, r.options.TTL).Err()

	return errors.Wrap(err, "setting key")
}

// SpecificationGetOne gets specification by id from cache.
func (r *Repository) SpecificationGetOne(ctxr context.Context, specificationID string) (specification.Specification, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.SpecificationGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr specification.Specification

	bytes, err := r.client.Get(ctx, specificationKey+specificationID).Bytes()
	if err != nil {
		return usr, errors.Wrap(err, "unable to get specification from cache")
	}

	if err = easyjson.Unmarshal(bytes, &usr); err != nil {
		return usr, errors.Wrap(err, "unable to unmarshal")
	}

	return usr, nil
}

// SpecificationCreate sets specification into cache.
func (r *Repository) SpecificationCreate(ctxr context.Context, usr specification.Specification) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.SpecificationCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, specificationKey+usr.ID, bytes, r.options.TTL).Err()

	return errors.Wrap(err, "setting key")
}

// SpecificationUpdate updates specification by id in cache.
func (r *Repository) SpecificationUpdate(ctxr context.Context, usr specification.Specification) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.SpecificationUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	r.client.Set(ctx, specificationKey+usr.ID, bytes, r.options.TTL)

	return nil
}

// SpecificationDelete deletes specification by id from cache.
func (r *Repository) SpecificationDelete(ctxr context.Context, specificationID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.SpecificationDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, specificationKey+specificationID).Err()

	return errors.Wrap(err, "deleting key")
}

// SpecificationInvalidate invalidate specifications cache.
func (r *Repository) SpecificationInvalidate(ctxr context.Context) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.SpecificationInvalidate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	iter := r.client.Scan(ctx, 0, specificationsKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return errors.Wrap(err, "deleting key")
		}
	}

	iter = r.client.Scan(ctx, 0, itemKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return errors.Wrap(err, "deleting key")
		}
	}

	iter = r.client.Scan(ctx, 0, itemsKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return errors.Wrap(err, "deleting key")
		}
	}

	return errors.Wrap(iter.Err(), "invalidate")
}
