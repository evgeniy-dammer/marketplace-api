package usecase

import "github.com/pkg/errors"

var (
	ErrInvalidHash          = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion  = errors.New("incompatible version of argon2")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidTokenClaims   = errors.New("token claims are not of type *tokenClaims")
	ErrUserNotFound         = errors.New("user not found")
	ErrUsersNotFound        = errors.New("users not found")
	ErrRolesNotFound        = errors.New("roles not found")
)
