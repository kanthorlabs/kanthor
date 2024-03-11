package api

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.SDK.API.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.SDK.API.ALREAD_STARTED")
)
