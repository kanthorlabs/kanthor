package scheduler

import (
	"context"
	"errors"
	"sync"

	"github.com/kanthorlabs/common/healthcheck"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/streaming"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/internal/constants"
	"github.com/kanthorlabs/kanthor/services/delivery/config"
	"github.com/kanthorlabs/kanthor/services/delivery/usecase"
)

type scheduler struct {
	conf   *config.Config
	logger logging.Logger
	infra  infrastructure.Infrastructure
	uc     usecase.Scheduler

	subscriber  streaming.Subscriber
	healthcheck healthcheck.Server

	mu     sync.Mutex
	status int
}

func (service *scheduler) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status == patterns.StatusStarted {
		return ErrAlreadyStarted
	}

	if err := service.infra.Connect(ctx); err != nil {
		return err
	}

	subscriber, err := service.infra.Streaming().Subscriber(constants.SubscriberMessage)
	if err != nil {
		return err
	}
	if err := subscriber.Connect(ctx); err != nil {
		return err
	}
	service.subscriber = subscriber

	service.status = patterns.StatusStarted
	service.logger.Info("started")
	return nil
}

func (service *scheduler) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status != patterns.StatusStarted {
		return ErrNotStarted
	}

	var returning error
	if err := service.subscriber.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.infra.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	service.status = patterns.StatusStopped
	service.logger.Info("stopped")

	return returning
}

func (service *scheduler) Run(ctx context.Context) error {
	topic := project.Subject(constants.TopicMessage)

	if err := service.subscriber.Sub(ctx, topic, handler(service)); err != nil {
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

func (service *scheduler) readiness() error {
	return service.healthcheck.Readiness(func() error {
		if err := service.subscriber.Readiness(); err != nil {
			return err
		}

		if err := service.infra.Readiness(); err != nil {
			return err
		}

		return nil
	})
}
