package redis

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/table"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

// TableGetAll gets tables from cache.
func (r *Repository) TableGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]table.Table, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.TableGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	tables := &table.ListTable{}

	bytes, err := r.client.Get(ctx, tablesKey+"o."+meta.OrganizationID).Bytes()
	if err != nil {
		return *tables, errors.Wrap(err, "unable to get tables from cache")
	}

	if err = easyjson.Unmarshal(bytes, tables); err != nil {
		return *tables, errors.Wrap(err, "unable to unmarshal")
	}

	return *tables, nil
}

// TableSetAll sets tables into cache.
func (r *Repository) TableSetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter, tables []table.Table) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.TableSetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	tableSlice := table.ListTable(tables)

	bytes, err := easyjson.Marshal(tableSlice)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, tablesKey+"o."+meta.OrganizationID, bytes, r.options.Ttl).Err()

	return err
}

// TableGetOne gets table by id from cache.
func (r *Repository) TableGetOne(ctxr context.Context, tableID string) (table.Table, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.TableGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr table.Table

	bytes, err := r.client.Get(ctx, tableKey+tableID).Bytes()
	if err != nil {
		return usr, errors.Wrap(err, "unable to get table from cache")
	}

	if err = easyjson.Unmarshal(bytes, &usr); err != nil {
		return usr, errors.Wrap(err, "unable to unmarshal")
	}

	return usr, nil
}

// TableCreate sets table into cache.
func (r *Repository) TableCreate(ctxr context.Context, usr table.Table) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.TableCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, tableKey+usr.ID, bytes, r.options.Ttl).Err()

	return err
}

// TableUpdate updates table by id in cache.
func (r *Repository) TableUpdate(ctxr context.Context, usr table.Table) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.TableUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	r.client.Set(ctx, tableKey+usr.ID, bytes, r.options.Ttl)

	return nil
}

// TableDelete deletes table by id from cache.
func (r *Repository) TableDelete(ctxr context.Context, tableID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.TableDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, tableKey+tableID).Err()

	return err
}

// TableInvalidate invalidate tables cache.
func (r *Repository) TableInvalidate(ctxr context.Context) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.TableInvalidate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	iter := r.client.Scan(ctx, 0, tablesKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}

	return iter.Err()
}
