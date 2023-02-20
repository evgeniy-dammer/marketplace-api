package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/pkg/errors"
)

// UserService is a user service.
type UserService struct {
	repo repository.User
}

// NewUserService is a constructor for UserService.
func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

// GetAll returns all users from the system.
func (s *UserService) GetAll(search string, status string, roleID string) ([]model.User, error) {
	users, err := s.repo.GetAll(search, status, roleID)

	return users, errors.Wrap(err, "users select failed")
}

// GetAllRoles returns all user roles from the system.
func (s *UserService) GetAllRoles() ([]model.Role, error) {
	roles, err := s.repo.GetAllRoles()

	return roles, errors.Wrap(err, "roles select failed")
}

// GetOne returns user by id from the system.
func (s *UserService) GetOne(userID string) (model.User, error) {
	user, err := s.repo.GetOne(userID)

	return user, errors.Wrap(err, "user select failed")
}

// Create hashes the password and insert User into system.
func (s *UserService) Create(userID string, user model.User) (string, error) {
	pass, err := generatePasswordHash(user.Password, params)
	if err != nil {
		return "", err
	}

	user.Password = pass

	ID, err := s.repo.Create(userID, user)

	return ID, errors.Wrap(err, "user create failed")
}

// Update updates user by id in the system.
func (s *UserService) Update(userID string, input model.UpdateUserInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	if input.Password != nil {
		pass, err := generatePasswordHash(*input.Password, params)
		if err != nil {
			return err
		}

		input.Password = &pass
	}

	return errors.Wrap(s.repo.Update(userID, input), "user update failed")
}

// Delete deletes user by id from the system.
func (s *UserService) Delete(userID string, dUserID string) error {
	err := s.repo.Delete(userID, dUserID)

	return errors.Wrap(err, "user delete failed")
}
