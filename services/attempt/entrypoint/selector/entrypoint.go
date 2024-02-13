package selector

import (
	"context"
	"errors"
	"sync"

	"github.com/kanthorlabs/common/healthcheck"
	"github.com/kanthorlabs/common/healthcheck/background"
	hcconfig "github.com/kanthorlabs/common/healthcheck/config"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/infrastructure/streaming"
	"github.com/kanthorlabs/kanthor/internal/constants"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/attempt/config"
	"github.com/kanthorlabs/kanthor/services/attempt/usecase"
	"github.com/kanthorlabs/kanthor/telemetry"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) (patterns.Runnable, error) {
	healthcheck, err := background.NewServer(hcconfig.Default("attempt.cronjob", 5000))
	if err != nil {
		return nil, err
	}

	logger = logger.With("service", "attempt", "entrypoint", "selector")
	return &selector{
		conf:       conf,
		logger:     logger,
		subscriber: infra.Stream.Subscriber("attempt_selector"),
		infra:      infra,
		db:         db,
		ds:         ds,
		uc:         uc,

		healthcheck: healthcheck,
	}, nil
}

type selector struct {
	conf       *config.Config
	logger     logging.Logger
	subscriber streaming.Subscriber
	infra      *infrastructure.Infrastructure
	db         database.Database
	ds         datastore.Datastore
	uc         usecase.Attempt

	healthcheck healthcheck.Server

	mu     sync.Mutex
	status int
}

func (service *selector) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status == patterns.StatusStarted {
		return ErrAlreadyStarted
	}

	if err := service.conf.Validate(); err != nil {
		return err
	}

	if err := service.db.Connect(ctx); err != nil {
		return err
	}

	if err := service.ds.Connect(ctx); err != nil {
		return err
	}

	if err := service.infra.Connect(ctx); err != nil {
		return err
	}

	if err := service.subscriber.Connect(ctx); err != nil {
		return err
	}

	if err := service.healthcheck.Connect(ctx); err != nil {
		return err
	}

	service.status = patterns.StatusStarted
	service.logger.Info("started")
	return nil
}

func (service *selector) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status != patterns.StatusStarted {
		return ErrNotStarted
	}
	service.status = patterns.StatusStopped
	service.logger.Info("stopped")

	var returning error
	if err := service.healthcheck.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.subscriber.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.infra.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.ds.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.db.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	return returning
}

func (service *selector) Run(ctx context.Context) error {
	tracectx := context.WithValue(ctx, telemetry.CtxTracer, telemetry.Tracer(project.Name("attempt_selector")))
	topic := constants.TopicAttemptTrigger
	if err := service.subscriber.Sub(tracectx, topic, Handler(service)); err != nil {
		return err
	}

	if err := service.readiness(); err != nil {
		return err
	}

	go func() {
		err := service.healthcheck.Liveness(func() error {
			if err := service.subscriber.Liveness(); err != nil {
				return err
			}

			if err := service.infra.Liveness(); err != nil {
				return err
			}

			if err := service.db.Liveness(); err != nil {
				return err
			}

			if err := service.ds.Liveness(); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			service.logger.Error(err)
		}
	}()

	service.logger.Infow("running", "topic", topic)
	<-ctx.Done()
	return nil
}

func (service *selector) readiness() error {
	return service.healthcheck.Readiness(func() error {
		if err := service.subscriber.Readiness(); err != nil {
			return err
		}

		if err := service.infra.Readiness(); err != nil {
			return err
		}

		if err := service.ds.Readiness(); err != nil {
			return err
		}

		if err := service.db.Readiness(); err != nil {
			return err
		}

		return nil
	})
}
