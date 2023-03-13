package authorization

import (
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// AuthorizationGetUserRole returns users role name.
func (s *UseCase) AuthorizationGetUserRole(ctx context.Context, userID string) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.AuthorizationGetUserRole")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getUserRoleWithCache(ctx, s, userID)
	}

	role, err := s.adapterStorage.AuthorizationGetUserRole(ctx, userID)

	return role, errors.Wrap(err, "can not get role")
}

func getUserRoleWithCache(ctx context.Context, s *UseCase, userID string) (string, error) {
	role, err := s.adapterCache.AuthorizationGetUserRole(ctx, userID)
	if err != nil {
		logger.Logger.Error("unable to get user role from cache", zap.String("error", err.Error()))
	}

	if role != "" {
		return role, nil
	}

	role, err = s.adapterStorage.AuthorizationGetUserRole(ctx, userID)

	if err != nil {
		return role, errors.Wrap(err, "user role select failed")
	}

	if err = s.adapterCache.AuthorizationSetUserRole(ctx, userID, role); err != nil {
		logger.Logger.Error("unable to add user role into cache", zap.String("error", err.Error()))
	}

	return role, nil
}
