package redis

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/organization"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

// OrganizationGetAll gets organizations from cache.
func (r *Repository) OrganizationGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]organization.Organization, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrganizationGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	organizations := &organization.ListOrganization{}

	bytes, err := r.client.Get(ctx, organizationsKey).Bytes()
	if err != nil {
		return *organizations, errors.Wrap(err, "unable to get organizations from cache")
	}

	if err = easyjson.Unmarshal(bytes, organizations); err != nil {
		return *organizations, errors.Wrap(err, "unable to unmarshal")
	}

	return *organizations, nil
}

// OrganizationSetAll sets organizations into cache.
func (r *Repository) OrganizationSetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter, organizations []organization.Organization) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrganizationSetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	organizationSlice := organization.ListOrganization(organizations)

	bytes, err := easyjson.Marshal(organizationSlice)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, organizationsKey, bytes, r.options.Ttl).Err()

	return err
}

// OrganizationGetOne gets organization by id from cache.
func (r *Repository) OrganizationGetOne(ctxr context.Context, organizationID string) (organization.Organization, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrganizationGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr organization.Organization

	bytes, err := r.client.Get(ctx, organizationKey+organizationID).Bytes()
	if err != nil {
		return usr, errors.Wrap(err, "unable to get organization from cache")
	}

	if err = easyjson.Unmarshal(bytes, &usr); err != nil {
		return usr, errors.Wrap(err, "unable to unmarshal")
	}

	return usr, nil
}

// OrganizationCreate sets organization into cache.
func (r *Repository) OrganizationCreate(ctxr context.Context, usr organization.Organization) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrganizationCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, organizationKey+usr.ID, bytes, r.options.Ttl).Err()

	return err
}

// OrganizationUpdate updates organization by id in cache.
func (r *Repository) OrganizationUpdate(ctxr context.Context, usr organization.Organization) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrganizationUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	r.client.Set(ctx, organizationKey+usr.ID, bytes, r.options.Ttl)

	return nil
}

// OrganizationDelete deletes organization by id from cache.
func (r *Repository) OrganizationDelete(ctxr context.Context, organizationID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrganizationDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, organizationKey+organizationID).Err()

	return err
}

// OrganizationInvalidate invalidate organizations cache.
func (r *Repository) OrganizationInvalidate(ctxr context.Context) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.OrganizationInvalidate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	iter := r.client.Scan(ctx, 0, organizationsKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}

	return iter.Err()
}
