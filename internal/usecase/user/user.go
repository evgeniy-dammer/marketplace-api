package user

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/role"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// UserGetAll returns all users from the system.
func (s *UseCase) UserGetAll(ctx context.Context, search string, status string, roleID string) ([]user.User, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.UserGetAll")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	users, err := s.adapterStorage.UserGetAll(ctx, search, status, roleID)

	return users, errors.Wrap(err, "users select failed")
}

// UserGetAllRoles returns all user roles from the system.
func (s *UseCase) UserGetAllRoles(ctx context.Context) ([]role.Role, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.UserGetAllRoles")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	roles, err := s.adapterStorage.UserGetAllRoles(ctx)

	return roles, errors.Wrap(err, "roles select failed")
}

// UserGetOne returns user by id from the system.
func (s *UseCase) UserGetOne(ctx context.Context, userID string) (user.User, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.UserGetOne")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var usr user.User

	if s.isCacheOn {
		return getOneWithCache(ctx, s, userID)
	}

	usr, err := s.adapterStorage.UserGetOne(ctx, userID)
	if err != nil {
		return usr, errors.Wrap(err, "user select failed")
	}

	return usr, nil
}

func getOneWithCache(ctx context.Context, s *UseCase, userID string) (user.User, error) {
	usr, err := s.adapterCache.UserGetOne(ctx, userID)
	if err != nil {
		logger.Logger.Error("unable to get user from cache", zap.String("error", err.Error()))
	}

	if usr != (user.User{}) {
		return usr, nil
	}

	usr, err = s.adapterStorage.UserGetOne(ctx, userID)

	if err != nil {
		return usr, errors.Wrap(err, "user select failed")
	}

	if err = s.adapterCache.UserCreate(ctx, usr); err != nil {
		logger.Logger.Error("unable to add user into cache", zap.String("error", err.Error()))
	}

	return usr, nil
}

// UserCreate hashes the password and insert User into system.
func (s *UseCase) UserCreate(ctx context.Context, userID string, input user.CreateUserInput) (string, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.UserCreate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	pass, err := usecase.GeneratePasswordHash(input.Password, usecase.Params)
	if err != nil {
		return "", err
	}

	input.Password = pass

	ID, err := s.adapterStorage.UserCreate(ctx, userID, input)
	if err != nil {
		return "", errors.Wrap(err, "user create in database failed")
	}

	if s.isCacheOn {
		usr, err := s.adapterStorage.UserGetOne(ctx, ID)
		if err != nil {
			return "", errors.Wrap(err, "user select from database failed")
		}

		err = s.adapterCache.UserCreate(ctx, usr)
		if err != nil {
			return "", errors.Wrap(err, "user create in cache failed")
		}
	}

	return ID, nil
}

// UserUpdate updates user by id in the system.
func (s *UseCase) UserUpdate(ctx context.Context, userID string, input user.UpdateUserInput) error {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.UserUpdate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	if input.Password != nil {
		pass, err := usecase.GeneratePasswordHash(*input.Password, usecase.Params)
		if err != nil {
			return err
		}

		input.Password = &pass
	}

	err := s.adapterStorage.UserUpdate(ctx, userID, input)
	if err != nil {
		return errors.Wrap(err, "user update in database failed")
	}

	if s.isCacheOn {
		usr, err := s.adapterStorage.UserGetOne(ctx, *input.ID)
		if err != nil {
			return errors.Wrap(err, "user select from database failed")
		}

		err = s.adapterCache.UserUpdate(ctx, usr)
		if err != nil {
			return errors.Wrap(err, "user update in cache failed")
		}
	}

	return nil
}

// UserDelete deletes user by id from the system.
func (s *UseCase) UserDelete(ctx context.Context, userID string, dUserID string) error {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.UserDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.UserDelete(ctx, userID, dUserID)
	if err != nil {
		return errors.Wrap(err, "user delete failed")
	}

	if s.isCacheOn {
		err = s.adapterCache.UserDelete(ctx, dUserID)
		if err != nil {
			return errors.Wrap(err, "user update in cache failed")
		}
	}

	return nil
}
