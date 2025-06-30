package repository

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrUserCreationFailed       = errors.New("failed to create user")
	ErrUserUpdateFailed         = errors.New("failed to update user")
	ErrUserDeletionFailed       = errors.New("failed to delete user")
	ErrUserRetrievalFailed      = errors.New("failed to retrieve user")
	ErrUserValidationFailed     = errors.New("failed to validate user")
	ErrUserAuthenticationFailed = errors.New("failed to authenticate user")
	ErrUserAuthorizationFailed  = errors.New("failed to authorize user")
)
