package domain

import "errors"

var (
	ErrNotFound       = errors.New("record not found")
	ErrConflict       = errors.New("record already exists")
	ErrInvalid        = errors.New("invalid input")
	ErrForbidden      = errors.New("operation not allowed")
	ErrAlreadyDone    = errors.New("task already completed")
	ErrDeviceInactive = errors.New("device is inactive")
)
