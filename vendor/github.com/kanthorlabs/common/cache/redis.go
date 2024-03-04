package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kanthorlabs/common/cache/config"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	goredis "github.com/redis/go-redis/v9"
)

// NewRedis creates a new redis cache instance that is using go-redis as the underlying engine.
func NewRedis(conf *config.Config, logger logging.Logger) (Cache, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	logger = logger.With("cache", "redis")
	return &redis{conf: conf, logger: logger}, nil
}

type redis struct {
	conf   *config.Config
	logger logging.Logger

	client *goredis.Client
	mu     sync.Mutex
	status int
}

func (instance *redis) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}
	conf, err := goredis.ParseURL(instance.conf.Uri)
	if err != nil {
		return err
	}
	instance.client = goredis.NewClient(conf)

	instance.status = patterns.StatusConnected
	return nil
}

func (instance *redis) Readiness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return instance.client.Ping(ctx).Err()
}

func (instance *redis) Liveness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return instance.client.Ping(ctx).Err()
}

func (instance *redis) Disconnect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	instance.status = patterns.StatusDisconnected

	var returning error
	if err := instance.client.Close(); err != nil {
		returning = errors.Join(returning, err)
	}
	instance.client = nil

	return returning
}

func (instance *redis) Get(ctx context.Context, key string, entry any) error {
	k, err := Key(key)
	if err != nil {
		return err
	}

	data, err := instance.client.Get(ctx, k).Bytes()
	if errors.Is(err, goredis.Nil) {
		return ErrEntryNotFound
	}

	return Unmarshal(data, entry)
}

func (instance *redis) Set(ctx context.Context, key string, entry any, ttl time.Duration) error {
	k, err := Key(key)
	if err != nil {
		return err
	}

	v, err := Marshal(entry)
	if err != nil {
		return fmt.Errorf("CACHE.VALUE.MARSHAL.ERROR: %w", err)
	}

	return instance.client.Set(ctx, k, v, ttl).Err()
}

func (instance *redis) Exist(ctx context.Context, key string) bool {
	k, err := Key(key)
	if err != nil {
		return false
	}

	entry, err := instance.client.Exists(ctx, k).Result()
	return err == nil && entry > 0
}

func (instance *redis) Del(ctx context.Context, key string) error {
	k, err := Key(key)
	if err != nil {
		return err
	}

	return instance.client.Del(ctx, k).Err()
}

func (instance *redis) Expire(ctx context.Context, key string, at time.Time) error {
	k, err := Key(key)
	if err != nil {
		return err
	}

	ttl := time.Until(at)
	if ttl < 0 {
		return errors.New("CACHE.TIME_TO_LIVE.NEGATIVE.ERROR")
	}

	ok, err := instance.client.ExpireAt(ctx, k, at).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ErrEntryNotFound
	}
	return nil
}
