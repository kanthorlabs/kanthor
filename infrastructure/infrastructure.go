package infrastructure

import (
	"context"

	"github.com/kanthorlabs/common/cache"
	"github.com/kanthorlabs/common/circuitbreaker"
	"github.com/kanthorlabs/common/distributedlockmanager"
	"github.com/kanthorlabs/common/idempotency"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/passport"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/common/streaming"
	"github.com/kanthorlabs/kanthor/infrastructure/config"
	"github.com/sourcegraph/conc/pool"
)

func New(conf *config.Config, logger logging.Logger) (Infrastructure, error) {
	logger = logger.With("component", "infrastructure")

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	infra := &infrastructure{conf: conf, logger: logger}
	var err error

	infra.database, err = database.New(&conf.Database, logger)
	if err != nil {
		return nil, err
	}
	infra.datastore, err = datastore.New(&conf.Datastore, logger)
	if err != nil {
		return nil, err
	}
	infra.streaming, err = streaming.New(&conf.Stream, logger)
	if err != nil {
		return nil, err
	}

	infra.cache, err = cache.New(&conf.Cache, logger)
	if err != nil {
		return nil, err
	}

	infra.distributedlockmanager, err = distributedlockmanager.New(&conf.DLM)
	if err != nil {
		return nil, err
	}

	infra.idempotency, err = idempotency.New(&conf.Idempotency, logger)
	if err != nil {
		return nil, err
	}

	infra.circuitbreaker, err = circuitbreaker.New(&conf.CircuitBreaker, logger)
	if err != nil {
		return nil, err
	}

	infra.passport, err = passport.New(&conf.Passport, logger)
	if err != nil {
		return nil, err
	}

	return infra, nil
}

type Infrastructure interface {
	patterns.Connectable
	Database() database.Database
	Datastore() datastore.Datastore
	Streaming() streaming.Stream
	Cache() cache.Cache
	DLM() distributedlockmanager.DistributedLockManager
	Idempotency() idempotency.Idempotency
	CircuitBreaker() circuitbreaker.CircuitBreaker
	Passport() passport.Passport
}

type infrastructure struct {
	conf   *config.Config
	logger logging.Logger

	database               database.Database
	datastore              datastore.Datastore
	streaming              streaming.Stream
	cache                  cache.Cache
	distributedlockmanager distributedlockmanager.DistributedLockManager
	idempotency            idempotency.Idempotency
	circuitbreaker         circuitbreaker.CircuitBreaker
	passport               passport.Passport
}

func (infra *infrastructure) Connect(ctx context.Context) error {
	p := pool.New().WithContext(ctx)

	p.Go(func(subctx context.Context) error {
		return infra.database.Connect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.datastore.Connect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.streaming.Connect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.cache.Connect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.distributedlockmanager.Connect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.idempotency.Connect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.passport.Connect(subctx)
	})

	if err := p.Wait(); err != nil {
		return err
	}

	infra.logger.Info("connected")
	return nil
}

func (infra *infrastructure) Readiness() error {
	p := pool.New().WithErrors()

	p.Go(func() error {
		return infra.database.Readiness()
	})
	p.Go(func() error {
		return infra.datastore.Readiness()
	})
	p.Go(func() error {
		return infra.streaming.Readiness()
	})
	p.Go(func() error {
		return infra.cache.Readiness()
	})
	p.Go(func() error {
		return infra.distributedlockmanager.Readiness()
	})
	p.Go(func() error {
		return infra.idempotency.Readiness()
	})
	p.Go(func() error {
		return infra.passport.Readiness()
	})

	if err := p.Wait(); err != nil {
		return err
	}

	infra.logger.Debug("ready")
	return nil
}

func (infra *infrastructure) Liveness() error {
	p := pool.New().WithErrors()

	p.Go(func() error {
		return infra.database.Readiness()
	})
	p.Go(func() error {
		return infra.datastore.Readiness()
	})
	p.Go(func() error {
		return infra.streaming.Readiness()
	})
	p.Go(func() error {
		return infra.cache.Readiness()
	})
	p.Go(func() error {
		return infra.distributedlockmanager.Readiness()
	})
	p.Go(func() error {
		return infra.idempotency.Readiness()
	})
	p.Go(func() error {
		return infra.passport.Readiness()
	})

	if err := p.Wait(); err != nil {
		return err
	}

	infra.logger.Debug("live")
	return nil
}

func (infra *infrastructure) Disconnect(ctx context.Context) error {
	p := pool.New().WithContext(ctx)

	p.Go(func(subctx context.Context) error {
		return infra.database.Disconnect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.datastore.Disconnect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.streaming.Disconnect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.cache.Disconnect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.distributedlockmanager.Disconnect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.idempotency.Disconnect(subctx)
	})
	p.Go(func(subctx context.Context) error {
		return infra.passport.Disconnect(subctx)
	})

	if err := p.Wait(); err != nil {
		return err
	}

	infra.logger.Debug("disconnected")
	return nil
}

func (infra *infrastructure) Database() database.Database {
	return infra.database
}

func (infra *infrastructure) Datastore() datastore.Datastore {
	return infra.datastore
}

func (infra *infrastructure) Streaming() streaming.Stream {
	return infra.streaming
}

func (infra *infrastructure) Cache() cache.Cache {
	return infra.cache
}

func (infra *infrastructure) DLM() distributedlockmanager.DistributedLockManager {
	return infra.distributedlockmanager
}

func (infra *infrastructure) Idempotency() idempotency.Idempotency {
	return infra.idempotency
}

func (infra *infrastructure) CircuitBreaker() circuitbreaker.CircuitBreaker {
	return infra.circuitbreaker
}

func (infra *infrastructure) Passport() passport.Passport {
	return infra.passport
}
