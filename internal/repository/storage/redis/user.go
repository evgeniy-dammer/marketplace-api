package redis

import (
	"encoding/json"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/role"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
)

// UserGetAll gets users from cache.
func (r *Repository) UserGetAll(ctxr context.Context, search string, status string, roleID string) ([]user.User, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.UserGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var users []user.User

	bytes, err := r.client.Get(ctx, usersKey+"s."+search+".t."+status+".r."+roleID).Bytes()
	if err != nil {
		return users, errors.Wrap(err, "unable to get users from cache")
	}

	if err = json.Unmarshal(bytes, &users); err != nil {
		return users, errors.Wrap(err, "unable to unmarshal")
	}

	return users, nil
}

// UserSetAll sets users into cache.
func (r *Repository) UserSetAll(ctxr context.Context, users []user.User, search string, status string, roleID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.UserSetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := json.Marshal(users)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, usersKey+"s."+search+".t."+status+".r."+roleID, bytes, r.options.Ttl).Err()

	return err
}

// UserGetOne gets user by id from cache.
func (r *Repository) UserGetOne(ctxr context.Context, userID string) (user.User, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.UserGetOne")
		defer span.End()

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
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.UserCreate")
		defer span.End()

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
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.UserUpdate")
		defer span.End()

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
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.UserDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, userKey+userID).Err()

	return err
}

// UserInvalidate invalidate users cache.
func (r *Repository) UserInvalidate(ctxr context.Context) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.UserInvalidate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	iter := r.client.Scan(ctx, 0, usersKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}

	return iter.Err()
}

// UserGetAllRoles gets all roles from cache.
func (r *Repository) UserGetAllRoles(ctxr context.Context) ([]role.Role, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.UserGetAllRoles")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var roles []role.Role

	bytes, err := r.client.Get(ctx, rolesKey).Bytes()
	if err != nil {
		return roles, errors.Wrap(err, "unable to get roles from cache")
	}

	if err = json.Unmarshal(bytes, &roles); err != nil {
		return roles, errors.Wrap(err, "unable to unmarshal")
	}

	return roles, nil
}

// UserSetAllRoles sets all roles in cache.
func (r *Repository) UserSetAllRoles(ctxr context.Context, roles []role.Role) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.UserSetAllRoles")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := json.Marshal(roles)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, rolesKey, bytes, r.options.Ttl).Err()

	return err
}
