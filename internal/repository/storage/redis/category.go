package redis

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/category"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

// CategoryGetAll gets categories from cache.
func (r *Repository) CategoryGetAll(ctxr context.Context, meta query.MetaData, _ queryparameter.QueryParameter) ([]category.Category, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CategoryGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	categories := &category.ListCategory{}

	bytes, err := r.client.Get(ctx, categoriesKey+"o."+meta.OrganizationID).Bytes()
	if err != nil {
		return *categories, errors.Wrap(err, "unable to get categories from cache")
	}

	if err = easyjson.Unmarshal(bytes, categories); err != nil {
		return *categories, errors.Wrap(err, "unable to unmarshal")
	}

	return *categories, nil
}

// CategorySetAll sets categories into cache.
func (r *Repository) CategorySetAll(ctxr context.Context, meta query.MetaData, _ queryparameter.QueryParameter, categories []category.Category) error { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CategorySetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	categorySlice := category.ListCategory(categories)

	bytes, err := easyjson.Marshal(categorySlice)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, categoriesKey+"o."+meta.OrganizationID, bytes, r.options.TTL).Err()

	return errors.Wrap(err, "setting key")
}

// CategoryGetOne gets category by id from cache.
func (r *Repository) CategoryGetOne(ctxr context.Context, categoryID string) (category.Category, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CategoryGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr category.Category

	bytes, err := r.client.Get(ctx, categoryKey+categoryID).Bytes()
	if err != nil {
		return usr, errors.Wrap(err, "unable to get category from cache")
	}

	if err = easyjson.Unmarshal(bytes, &usr); err != nil {
		return usr, errors.Wrap(err, "unable to unmarshal")
	}

	return usr, nil
}

// CategoryCreate sets category into cache.
func (r *Repository) CategoryCreate(ctxr context.Context, usr category.Category) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CategoryCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, categoryKey+usr.ID, bytes, r.options.TTL).Err()

	return errors.Wrap(err, "setting key")
}

// CategoryUpdate updates category by id in cache.
func (r *Repository) CategoryUpdate(ctxr context.Context, usr category.Category) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CategoryUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	r.client.Set(ctx, categoryKey+usr.ID, bytes, r.options.TTL)

	return nil
}

// CategoryDelete deletes category by id from cache.
func (r *Repository) CategoryDelete(ctxr context.Context, categoryID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CategoryDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, categoryKey+categoryID).Err()

	return errors.Wrap(err, "deleting key")
}

// CategoryInvalidate invalidate categories cache.
func (r *Repository) CategoryInvalidate(ctxr context.Context) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CategoryInvalidate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	iter := r.client.Scan(ctx, 0, categoriesKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return errors.Wrap(err, "deleting key")
		}
	}

	return errors.Wrap(iter.Err(), "invalidate")
}
