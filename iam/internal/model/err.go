package model

import "errors"

var (
	ErrBadRequest = errors.New("bad request")
	ErrInternal   = errors.New("internal server error")

	ErrInvalidCredentials = errors.New("invalid login or password")

	ErrUserNotFound            = errors.New("user not found")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrFailedToCreateUser      = errors.New("failed to create user")
	ErrFailedToUpdateUser      = errors.New("failed to update user")
	ErrFailedToDeleteUser      = errors.New("failed to delete user")
	ErrFailedToGetUser         = errors.New("failed to get user")
	ErrUserConstraintViolation = errors.New("user constraint violation")
	ErrInvalidUserData         = errors.New("invalid user data")

	ErrSessionNotFound       = errors.New("session not found")
	ErrSessionExpired        = errors.New("session expired")
	ErrFailedToCreateSession = errors.New("failed to create session")
	ErrFailedToDeleteSession = errors.New("failed to delete session")
	ErrFailedToStoreInCache  = errors.New("failed to store in cache")
	ErrFailedToReadFromCache = errors.New("failed to read from cache")
	ErrInvalidSessionData    = errors.New("invalid session data")

	ErrNotificationNotFound      = errors.New("notification method not found")
	ErrNotificationAlreadyExists = errors.New("notification method already exists")

	ErrFailedToCreateNotification = errors.New("failed to create notification method")
	ErrFailedToDeleteNotification = errors.New("failed to delete notification method")
	ErrFailedToGetNotification    = errors.New("failed to get notification method")
	ErrFailedToListNotifications  = errors.New("failed to list notification methods")

	ErrInvalidNotificationData = errors.New("invalid notification data")

	ErrNotificationUserConstraintViolation = errors.New("notification user constraint violation")
)
