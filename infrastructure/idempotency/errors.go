package idempotency

import "errors"

var (
	ErrNotConnected     = errors.New("IDEMPOTENCY.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("IDEMPOTENCY.ALREADY_CONNECTED.ERROR")
)
