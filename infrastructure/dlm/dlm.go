package dlm

import (
	"context"
	"fmt"
)

type Factory func(key string, opts ...Option) DistributedLockManager

type DistributedLockManager interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}

func New(conf *Config) (Factory, error) {
	return NewRedlock(conf)
}

func Key(key string) string {
	return fmt.Sprintf("dlm/%s", key)
}
