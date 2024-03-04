package api

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.PORTAL.API.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.PORTAL.API.ALREAD_STARTED")
)
