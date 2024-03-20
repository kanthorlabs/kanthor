package dispatcher

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.DELIVERY.DISPATCHER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.DELIVERY.DISPATCHER.ALREAD_STARTED")
)
