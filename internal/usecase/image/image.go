package image

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/image"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// ImageGetAll returns all images from the system.
func (s *UseCase) ImageGetAll(ctx context.Context, userID string, organizationID string) ([]image.Image, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.ImageGetAll")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	images, err := s.adapterStorage.ImageGetAll(ctx, userID, organizationID)

	return images, errors.Wrap(err, "images select error")
}

// ImageGetOne returns image by id from the system.
func (s *UseCase) ImageGetOne(ctx context.Context, userID string, organizationID string, imageID string) (image.Image, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.ImageGetOne")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	imageSingle, err := s.adapterStorage.ImageGetOne(ctx, userID, organizationID, imageID)

	return imageSingle, errors.Wrap(err, "image select error")
}

// ImageCreate inserts image into system.
func (s *UseCase) ImageCreate(ctx context.Context, userID string, input image.CreateImageInput) (string, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.ImageCreate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	imageID, err := s.adapterStorage.ImageCreate(ctx, userID, input)

	return imageID, errors.Wrap(err, "image create error")
}

// ImageUpdate updates image by id in the system.
func (s *UseCase) ImageUpdate(ctx context.Context, userID string, input image.UpdateImageInput) error {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.ImageUpdate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.ImageUpdate(ctx, userID, input)

	return errors.Wrap(err, "image update error")
}

// ImageDelete deletes image by id from the system.
func (s *UseCase) ImageDelete(ctx context.Context, userID string, organizationID string, imageID string) error {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.ImageDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.ImageDelete(ctx, userID, organizationID, imageID)

	return errors.Wrap(err, "image delete error")
}
