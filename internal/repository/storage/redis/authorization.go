package redis

import (
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// AuthorizationGetUserRole gets users role name from cache.
func (r *Repository) AuthorizationGetUserRole(ctxr context.Context, userID string) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.AuthorizationGetUserRole")
		defer span.End()

		ctx = context.New(ctxt)
	}

	rle, err := r.client.Get(ctx, roleKey+userID).Result()
	if err != nil {
		return "", errors.Wrap(err, "getting key")
	}

	return rle, nil
}

// AuthorizationSetUserRole sets user role into cache.
func (r *Repository) AuthorizationSetUserRole(ctxr context.Context, userID string, rle string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.Authorization	SetUserRole")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Set(ctx, roleKey+userID, rle, r.options.TTL).Err()

	return errors.Wrap(err, "setting key")
}
