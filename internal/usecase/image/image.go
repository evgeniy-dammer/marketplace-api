package image

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/image"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ImageGetAll returns all images from the system.
func (s *UseCase) ImageGetAll(ctx context.Context, userID string, organizationID string) ([]image.Image, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.ImageGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getAllWithCache(ctx, s, userID, organizationID)
	}

	images, err := s.adapterStorage.ImageGetAll(ctx, userID, organizationID)

	return images, errors.Wrap(err, "images select error")
}

// getAllWithCache returns images from cache if exists.
func getAllWithCache(ctx context.Context, s *UseCase, userID string, organizationID string) ([]image.Image, error) {
	images, err := s.adapterCache.ImageGetAll(ctx, organizationID)
	if err != nil {
		logger.Logger.Error("unable to get images from cache", zap.String("error", err.Error()))
	}

	if len(images) > 0 {
		return images, nil
	}

	images, err = s.adapterStorage.ImageGetAll(ctx, userID, organizationID)
	if err != nil {
		return images, errors.Wrap(err, "images select failed")
	}

	if err = s.adapterCache.ImageSetAll(ctx, organizationID, images); err != nil {
		logger.Logger.Error("unable to add images into cache", zap.String("error", err.Error()))
	}

	return images, nil
}

// ImageGetOne returns image by id from the system.
func (s *UseCase) ImageGetOne(ctx context.Context, userID string, organizationID string, imageID string) (image.Image, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.ImageGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getOneWithCache(ctx, s, userID, organizationID, imageID)
	}

	imageSingle, err := s.adapterStorage.ImageGetOne(ctx, userID, organizationID, imageID)

	return imageSingle, errors.Wrap(err, "image select error")
}

// getOneWithCache returns image by id from cache if exists.
func getOneWithCache(ctx context.Context, s *UseCase, userID string, organizationID string, imageID string) (image.Image, error) {
	img, err := s.adapterCache.ImageGetOne(ctx, imageID)
	if err != nil {
		logger.Logger.Error("unable to get image from cache", zap.String("error", err.Error()))
	}

	if img != (image.Image{}) {
		return img, nil
	}

	img, err = s.adapterStorage.ImageGetOne(ctx, userID, organizationID, imageID)
	if err != nil {
		return img, errors.Wrap(err, "image select failed")
	}

	if err = s.adapterCache.ImageCreate(ctx, img); err != nil {
		logger.Logger.Error("unable to add image into cache", zap.String("error", err.Error()))
	}

	return img, nil
}

// ImageCreate inserts image into system.
func (s *UseCase) ImageCreate(ctx context.Context, userID string, input image.CreateImageInput) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.ImageCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	imageID, err := s.adapterStorage.ImageCreate(ctx, userID, input)
	if err != nil {
		return imageID, errors.Wrap(err, "image create error")
	}

	if s.isCacheOn {
		img, err := s.adapterStorage.ImageGetOne(ctx, userID, input.OrganizationID, imageID)
		if err != nil {
			return "", errors.Wrap(err, "image select from database failed")
		}

		err = s.adapterCache.ImageCreate(ctx, img)
		if err != nil {
			return "", errors.Wrap(err, "image create in cache failed")
		}

		err = s.adapterCache.ImageInvalidate(ctx)
		if err != nil {
			return "", errors.Wrap(err, "image invalidate users in cache failed")
		}
	}

	return imageID, nil
}

// ImageUpdate updates image by id in the system.
func (s *UseCase) ImageUpdate(ctx context.Context, userID string, input image.UpdateImageInput) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.ImageUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.ImageUpdate(ctx, userID, input)
	if err != nil {
		return errors.Wrap(err, "image update in database failed")
	}

	if s.isCacheOn {
		img, err := s.adapterStorage.ImageGetOne(ctx, userID, *input.OrganizationID, *input.ID)
		if err != nil {
			return errors.Wrap(err, "image select from database failed")
		}

		err = s.adapterCache.ImageUpdate(ctx, img)
		if err != nil {
			return errors.Wrap(err, "image update in cache failed")
		}

		err = s.adapterCache.ImageInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "image invalidate users in cache failed")
		}
	}

	return nil
}

// ImageDelete deletes image by id from the system.
func (s *UseCase) ImageDelete(ctx context.Context, userID string, organizationID string, imageID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.ImageDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.ImageDelete(ctx, userID, organizationID, imageID)
	if err != nil {
		return errors.Wrap(err, "image delete failed")
	}

	if s.isCacheOn {
		err = s.adapterCache.ImageDelete(ctx, imageID)
		if err != nil {
			return errors.Wrap(err, "image update in cache failed")
		}

		err = s.adapterCache.ImageInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "invalidate images in cache failed")
		}
	}

	return nil
}
