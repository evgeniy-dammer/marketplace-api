package redis

import (
	"encoding/json"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// UserGetOne gets user by id from cache.
func (r *Repository) UserGetOne(ctxr context.Context, userID string) (user.User, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "Cache.UserGetOne")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var usr user.User

	bytes, err := r.client.Get(ctx, userKey+userID).Bytes()
	if err != nil {
		return usr, errors.Wrap(err, "unable to get user from cache")
	}

	if err = json.Unmarshal(bytes, &usr); err != nil {
		return usr, errors.Wrap(err, "unable to unmarshal")
	}

	return usr, nil
}

// UserCreate sets user into cache.
func (r *Repository) UserCreate(ctxr context.Context, usr user.User) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "Cache.UserCreate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	bytes, err := json.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, userKey+usr.ID, bytes, r.options.Ttl).Err()

	return err
}

// UserUpdate updates user by id in cache.
func (r *Repository) UserUpdate(ctxr context.Context, usr user.User) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "Cache.UserUpdate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	bytes, err := json.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	r.client.Set(ctx, userKey+usr.ID, bytes, r.options.Ttl)

	return nil
}

// UserDelete deletes user by id from cache.
func (r *Repository) UserDelete(ctxr context.Context, userID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "Cache.UserDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, userKey+userID).Err()

	return err
}
