package redis

import (
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
)

// AuthenticationGetUserRole gets users role name from cache
func (r *Repository) AuthenticationGetUserRole(ctxr context.Context, userID string) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "Cache.AuthenticationGetUserRole")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	rle, err := r.client.Get(ctx, roleKey+userID).Result()
	if err != nil {
		return "", err
	}

	return rle, nil
}

// AuthenticationSetUserRole sets user role into cache
func (r *Repository) AuthenticationSetUserRole(ctxr context.Context, userID string, rle string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "Cache.AuthenticationSetUserRole")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	err := r.client.Set(ctx, roleKey+userID, rle, r.options.Ttl).Err()

	return err
}
