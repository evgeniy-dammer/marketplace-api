package user

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/role"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// UserGetAll returns all users from the system.
func (s *UseCase) UserGetAll(ctx context.Context, search string, status string, roleID string) ([]user.User, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.UserGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getAllWithCache(ctx, s, search, status, roleID)
	}

	users, err := s.adapterStorage.UserGetAll(ctx, search, status, roleID)

	return users, errors.Wrap(err, "users select failed")
}

// getAllWithCache returns users from cache if exists.
func getAllWithCache(ctx context.Context, s *UseCase, search string, status string, roleID string) ([]user.User, error) {
	users, err := s.adapterCache.UserGetAll(ctx, search, status, roleID)
	if err != nil {
		logger.Logger.Error("unable to get users from cache", zap.String("error", err.Error()))
	}

	if len(users) > 0 {
		return users, nil
	}

	users, err = s.adapterStorage.UserGetAll(ctx, search, status, roleID)

	if err != nil {
		return users, errors.Wrap(err, "users select failed")
	}

	if err = s.adapterCache.UserSetAll(ctx, users, search, status, roleID); err != nil {
		logger.Logger.Error("unable to add users into cache", zap.String("error", err.Error()))
	}

	return users, nil
}

// UserGetAllRoles returns all user roles from the system.
func (s *UseCase) UserGetAllRoles(ctx context.Context) ([]role.Role, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.UserGetAllRoles")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getAllRolesWithCache(ctx, s)
	}

	roles, err := s.adapterStorage.UserGetAllRoles(ctx)

	return roles, errors.Wrap(err, "roles select failed")
}

// getAllRolesWithCache returns all user roles from cache if exists.
func getAllRolesWithCache(ctx context.Context, s *UseCase) ([]role.Role, error) {
	roles, err := s.adapterCache.UserGetAllRoles(ctx)
	if err != nil {
		logger.Logger.Error("unable to get roles from cache", zap.String("error", err.Error()))
	}

	if len(roles) > 0 {
		return roles, nil
	}

	roles, err = s.adapterStorage.UserGetAllRoles(ctx)

	if err != nil {
		return roles, errors.Wrap(err, "roles select failed")
	}

	if err = s.adapterCache.UserSetAllRoles(ctx, roles); err != nil {
		logger.Logger.Error("unable to add roles into cache", zap.String("error", err.Error()))
	}

	return roles, nil
}

// UserGetOne returns user by id from the system.
func (s *UseCase) UserGetOne(ctx context.Context, userID string) (user.User, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.UserGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getOneWithCache(ctx, s, userID)
	}

	usr, err := s.adapterStorage.UserGetOne(ctx, userID)
	if err != nil {
		return usr, errors.Wrap(err, "user select failed")
	}

	return usr, nil
}

// getOneWithCache returns user by id from cache if exists.
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
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.UserCreate")
		defer span.End()

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

		err = s.adapterCache.UserInvalidate(ctx)
		if err != nil {
			return "", errors.Wrap(err, "user invalidate users in cache failed")
		}
	}

	return ID, nil
}

// UserUpdate updates user by id in the system.
func (s *UseCase) UserUpdate(ctx context.Context, userID string, input user.UpdateUserInput) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.UserUpdate")
		defer span.End()

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

		err = s.adapterCache.UserInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "user invalidate users in cache failed")
		}
	}

	return nil
}

// UserDelete deletes user by id from the system.
func (s *UseCase) UserDelete(ctx context.Context, userID string, dUserID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.UserDelete")
		defer span.End()

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

		err = s.adapterCache.UserInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "invalidate users in cache failed")
		}
	}

	return nil
}
