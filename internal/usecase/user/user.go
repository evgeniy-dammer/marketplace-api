package user

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/role"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// UserGetAll returns all users from the system.
func (s *UseCase) UserGetAll(search string, status string, roleID string) ([]user.User, error) {
	users, err := s.adapterStorage.UserGetAll(search, status, roleID)

	return users, errors.Wrap(err, "users select failed")
}

// UserGetAllRoles returns all user roles from the system.
func (s *UseCase) UserGetAllRoles() ([]role.Role, error) {
	roles, err := s.adapterStorage.UserGetAllRoles()

	return roles, errors.Wrap(err, "roles select failed")
}

// UserGetOne returns user by id from the system.
func (s *UseCase) UserGetOne(userID string) (user.User, error) {
	var usr user.User

	if s.adapterCache != nil {
		return getOneWithCache(s, userID)
	}

	usr, err := s.adapterStorage.UserGetOne(userID)
	if err != nil {
		return usr, errors.Wrap(err, "user select failed")
	}

	return usr, nil
}

func getOneWithCache(s *UseCase, userID string) (user.User, error) {
	usr, err := s.adapterCache.UserGetOne(userID)
	if err != nil {
		logger.Logger.Error("unable to get user from cache", zap.String("error", err.Error()))
	}

	if usr != (user.User{}) {
		return usr, nil
	}

	usr, err = s.adapterStorage.UserGetOne(userID)

	if err != nil {
		return usr, errors.Wrap(err, "user select failed")
	}

	if err = s.adapterCache.UserCreate(userID, usr); err != nil {
		logger.Logger.Error("unable to add user into cache", zap.String("error", err.Error()))
	}

	return usr, nil
}

// UserCreate hashes the password and insert User into system.
func (s *UseCase) UserCreate(userID string, user user.User) (string, error) {
	pass, err := usecase.GeneratePasswordHash(user.Password, usecase.Params)
	if err != nil {
		return "", err
	}

	user.Password = pass

	ID, err := s.adapterStorage.UserCreate(userID, user)

	return ID, errors.Wrap(err, "user create failed")
}

// UserUpdate updates user by id in the system.
func (s *UseCase) UserUpdate(userID string, input user.UpdateUserInput) error {
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

	return errors.Wrap(s.adapterStorage.UserUpdate(userID, input), "user update failed")
}

// UserDelete deletes user by id from the system.
func (s *UseCase) UserDelete(userID string, dUserID string) error {
	err := s.adapterStorage.UserDelete(userID, dUserID)

	return errors.Wrap(err, "user delete failed")
}
