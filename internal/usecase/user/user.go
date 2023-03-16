package user

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/role"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/user"
	"github.com/evgeniy-dammer/marketplace-api/internal/usecase"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// UserGetAll returns all users from the system.
func (s *UseCase) UserGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]user.User, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.UserGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getAllWithCache(ctx, meta, params)
	}

	users, err := s.adapterStorage.UserGetAll(ctx, meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "users select failed")
	}

	return users, nil
}

// getAllWithCache returns users from cache if exists.
func (s *UseCase) getAllWithCache(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]user.User, error) {
	users, err := s.adapterCache.UserGetAll(ctx, meta, params)
	if err != nil {
		logger.Logger.Error("unable to get users from cache", zap.String("error", err.Error()))
	}

	if len(users) > 0 {
		return users, nil
	}

	users, err = s.adapterStorage.UserGetAll(ctx, meta, params)
	if err != nil {
		return users, errors.Wrap(err, "users select failed")
	}

	if err = s.adapterCache.UserSetAll(ctx, meta, params, users); err != nil {
		logger.Logger.Error("unable to add users into cache", zap.String("error", err.Error()))
	}

	return users, nil
}

// UserGetAllRoles returns all user roles from the system.
func (s *UseCase) UserGetAllRoles(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]role.Role, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.UserGetAllRoles")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getAllRolesWithCache(ctx, meta, params)
	}

	roles, err := s.adapterStorage.UserGetAllRoles(ctx, meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "roles select failed")
	}

	return roles, nil
}

// getAllRolesWithCache returns all user roles from cache if exists.
func (s *UseCase) getAllRolesWithCache(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]role.Role, error) {
	roles, err := s.adapterCache.UserGetAllRoles(ctx, meta, params)
	if err != nil {
		logger.Logger.Error("unable to get roles from cache", zap.String("error", err.Error()))
	}

	if len(roles) > 0 {
		return roles, nil
	}

	roles, err = s.adapterStorage.UserGetAllRoles(ctx, meta, params)
	if err != nil {
		return roles, errors.Wrap(err, "roles select failed")
	}

	if err = s.adapterCache.UserSetAllRoles(ctx, meta, params, roles); err != nil {
		logger.Logger.Error("unable to add roles into cache", zap.String("error", err.Error()))
	}

	return roles, nil
}

// UserGetOne returns user by id from the system.
func (s *UseCase) UserGetOne(ctx context.Context, meta query.MetaData, userID string) (user.User, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.UserGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getOneWithCache(ctx, meta, userID)
	}

	usr, err := s.adapterStorage.UserGetOne(ctx, meta, userID)
	if err != nil {
		return user.User{}, errors.Wrap(err, "user select failed")
	}

	return usr, nil
}

// getOneWithCache returns user by id from cache if exists.
func (s *UseCase) getOneWithCache(ctx context.Context, meta query.MetaData, userID string) (user.User, error) {
	usr, err := s.adapterCache.UserGetOne(ctx, userID)
	if err != nil {
		logger.Logger.Error("unable to get user from cache", zap.String("error", err.Error()))
	}

	if usr != (user.User{}) {
		return usr, nil
	}

	usr, err = s.adapterStorage.UserGetOne(ctx, meta, userID)
	if err != nil {
		return usr, errors.Wrap(err, "user select failed")
	}

	if err = s.adapterCache.UserCreate(ctx, usr); err != nil {
		logger.Logger.Error("unable to add user into cache", zap.String("error", err.Error()))
	}

	return usr, nil
}

// UserCreate hashes the password and insert User into system.
func (s *UseCase) UserCreate(ctx context.Context, meta query.MetaData, input user.CreateUserInput) (string, error) {
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

	ID, err := s.adapterStorage.UserCreate(ctx, meta, input)
	if err != nil {
		return "", errors.Wrap(err, "user create in database failed")
	}

	if s.isCacheOn {
		usr, err := s.adapterStorage.UserGetOne(ctx, meta, ID)
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
func (s *UseCase) UserUpdate(ctx context.Context, meta query.MetaData, input user.UpdateUserInput) error {
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

	err := s.adapterStorage.UserUpdate(ctx, meta, input)
	if err != nil {
		return errors.Wrap(err, "user update in database failed")
	}

	if s.isCacheOn {
		usr, err := s.adapterStorage.UserGetOne(ctx, meta, *input.ID)
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
func (s *UseCase) UserDelete(ctx context.Context, meta query.MetaData, dUserID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.UserDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.UserDelete(ctx, meta, dUserID)
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
