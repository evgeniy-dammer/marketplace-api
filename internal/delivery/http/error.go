package http

import "errors"

var (
	ErrEmptyIDParam      = errors.New("empty id param")
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrEmptyAuthHeader   = errors.New("empty auth header")
	ErrUserIsNotFound    = errors.New("user is not found")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrRoleIsNotFound    = errors.New("role is not found")
	ErrAccessDenied      = errors.New("access denied")
)
