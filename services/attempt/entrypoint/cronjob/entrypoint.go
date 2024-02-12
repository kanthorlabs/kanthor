package cronjob

import (
	"context"
	"errors"
	"sync"

	"github.com/kanthorlabs/common/healthcheck"
	"github.com/kanthorlabs/common/healthcheck/background"
	hcconfig "github.com/kanthorlabs/common/healthcheck/config"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/kanthor/database"
	"github.com/kanthorlabs/kanthor/datastore"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/attempt/config"
	"github.com/kanthorlabs/kanthor/services/attempt/usecase"
	"github.com/robfig/cron/v3"
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

	logger = logger.With("service", "attempt", "entrypoint", "cronjob")
	return &cronjob{
		conf:   conf,
		logger: logger,
		infra:  infra,
		db:     db,
		ds:     ds,
		uc:     uc,

		cron:        cron.New(),
		healthcheck: healthcheck,
	}, nil
}

type cronjob struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	db     database.Database
	ds     datastore.Datastore
	uc     usecase.Attempt

	cron        *cron.Cron
	healthcheck healthcheck.Server

	mu     sync.Mutex
	status int
}

func (service *cronjob) Start(ctx context.Context) error {
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

	if err := service.healthcheck.Connect(ctx); err != nil {
		return err
	}

	_, err := service.cron.AddFunc(service.conf.Cronjob.Scheduler, UseJob(service))
	if err != nil {
		return err
	}

	service.status = patterns.StatusStarted
	service.logger.Info("started")
	return nil
}

func (service *cronjob) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status != patterns.StatusStarted {
		return ErrNotStarted
	}
	service.status = patterns.StatusStopped
	service.logger.Info("stopped")

	// wait for the cronjob is done
	<-service.cron.Stop().Done()

	var returning error
	if err := service.healthcheck.Disconnect(ctx); err != nil {
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

func (service *cronjob) Run(ctx context.Context) error {
	if err := service.readiness(); err != nil {
		return err
	}

	go func() {
		err := service.healthcheck.Liveness(func() error {
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

	// in development enviroment, we want to run all jobs after the startup process
	if project.IsDev() {
		entries := service.cron.Entries()
		for _, entry := range entries {
			entry.Job.Run()
		}
	}

	service.logger.Infow("running")
	go func() {
		service.cron.Run()
	}()
	<-ctx.Done()
	return nil
}

func (service *cronjob) readiness() error {
	return service.healthcheck.Readiness(func() error {
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
