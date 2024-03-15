package scheduler

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.DELIVERY.SCHEDULER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.DELIVERY.SCHEDULER.ALREAD_STARTED")
)
