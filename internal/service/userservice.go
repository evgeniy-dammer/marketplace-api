package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
)

// UserService is a user service
type UserService struct {
	repo repository.User
}

// NewUserService is a constructor for UserService
func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

// GetAll returns all users from the system
func (s *UserService) GetAll(search string, status string, roleId string) ([]model.User, error) {
	return s.repo.GetAll(search, status, roleId)
}

// GetAllRoles returns all user roles from the system
func (s *UserService) GetAllRoles() ([]model.Role, error) {
	return s.repo.GetAllRoles()
}

// GetOne returns user by id from the system
func (s *UserService) GetOne(userId string) (model.User, error) {
	return s.repo.GetOne(userId)
}

// Create hashes the password and insert User into system
func (s *UserService) Create(user model.User, statusId string) (string, error) {
	pass, err := generatePasswordHash(user.Password, params)

	if err != nil {
		return "", err
	}

	user.Password = pass

	return s.repo.Create(user, statusId)
}

// Update updates user by id in the system
func (s *UserService) Update(userId string, input model.UpdateUserInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	if input.Password != nil {
		pass, err := generatePasswordHash(*input.Password, params)

		if err != nil {
			return err
		}

		input.Password = &pass
	}

	return s.repo.Update(userId, input)
}

// Delete deletes user by id from the system
func (s *UserService) Delete(userId string) error {
	return s.repo.Delete(userId)
}

// GetActiveStatusId returns ID for active status from the system
func (s *UserService) GetActiveStatusId(name string) (string, error) {
	return s.repo.GetActiveStatusId(name)
}
