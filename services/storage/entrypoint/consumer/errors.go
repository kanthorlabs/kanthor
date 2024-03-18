package consumer

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.STORAGE.CONSUMER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.STORAGE.CONSUMER.ALREAD_STARTED")
)
