package redis

import (
	"encoding/json"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/image"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
)

// ImageGetAll gets images from cache.
func (r *Repository) ImageGetAll(ctxr context.Context, organizationID string) ([]image.Image, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ImageGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var images []image.Image

	bytes, err := r.client.Get(ctx, imagesKey+"o."+organizationID).Bytes()
	if err != nil {
		return images, errors.Wrap(err, "unable to get images from cache")
	}

	if err = json.Unmarshal(bytes, &images); err != nil {
		return images, errors.Wrap(err, "unable to unmarshal")
	}

	return images, nil
}

// ImageSetAll sets images into cache.
func (r *Repository) ImageSetAll(ctxr context.Context, organizationID string, images []image.Image) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ImageSetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := json.Marshal(images)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, imagesKey+"o."+organizationID, bytes, r.options.Ttl).Err()

	return err
}

// ImageGetOne gets image by id from cache.
func (r *Repository) ImageGetOne(ctxr context.Context, imageID string) (image.Image, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ImageGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr image.Image

	bytes, err := r.client.Get(ctx, imageKey+imageID).Bytes()
	if err != nil {
		return usr, errors.Wrap(err, "unable to get image from cache")
	}

	if err = json.Unmarshal(bytes, &usr); err != nil {
		return usr, errors.Wrap(err, "unable to unmarshal")
	}

	return usr, nil
}

// ImageCreate sets image into cache.
func (r *Repository) ImageCreate(ctxr context.Context, usr image.Image) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ImageCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := json.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, imageKey+usr.ID, bytes, r.options.Ttl).Err()

	return err
}

// ImageUpdate updates image by id in cache.
func (r *Repository) ImageUpdate(ctxr context.Context, usr image.Image) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ImageUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := json.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	r.client.Set(ctx, imageKey+usr.ID, bytes, r.options.Ttl)

	return nil
}

// ImageDelete deletes image by id from cache.
func (r *Repository) ImageDelete(ctxr context.Context, imageID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ImageDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, imageKey+imageID).Err()

	return err
}

// ImageInvalidate invalidate images cache.
func (r *Repository) ImageInvalidate(ctxr context.Context) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.ImageInvalidate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	iter := r.client.Scan(ctx, 0, imagesKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}

	return iter.Err()
}
