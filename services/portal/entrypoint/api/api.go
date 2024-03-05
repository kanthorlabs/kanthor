package api

import (
	"context"
	"errors"
	"sync"

	"github.com/kanthorlabs/common/gateway"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
	"github.com/sourcegraph/conc/pool"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra infrastructure.Infrastructure,
	uc usecase.Portal,
) (patterns.Runnable, error) {
	entrypoint := &portal{
		conf:   conf,
		logger: logger.With("entrypoint", "api"),
		infra:  infra,
		uc:     uc,
	}
	return entrypoint, nil
}

type portal struct {
	conf   *config.Config
	logger logging.Logger
	infra  infrastructure.Infrastructure
	uc     usecase.Portal

	server gateway.Gateway
	mu     sync.Mutex
	status int
}

func (service *portal) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status == patterns.StatusStarted {
		return ErrAlreadyStarted
	}

	p := pool.New().WithContext(ctx)
	p.Go(func(subctx context.Context) error {
		return service.infra.Connect(subctx)
	})
	if err := p.Wait(); err != nil {
		return err
	}

	server, err := gateway.New(&service.conf.Portal.Gateway, service.logger)
	if err != nil {
		return err
	}
	service.server = server
	if err := service.server.Start(ctx); err != nil {
		return err
	}

	service.status = patterns.StatusStarted
	service.logger.Info("started")
	return nil
}

func (service *portal) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status != patterns.StatusStarted {
		return ErrNotStarted
	}
	service.status = patterns.StatusStopped

	var returning error
	if err := service.server.Stop(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	p := pool.New().WithContext(ctx)
	p.Go(func(subctx context.Context) error {
		return service.infra.Disconnect(subctx)
	})
	if err := p.Wait(); err != nil {
		returning = errors.Join(returning, err)
	}

	service.logger.Info("stopped")
	return returning
}

func (service *portal) Run(ctx context.Context) error {
	if err := service.httpx(); err != nil {
		return err
	}

	service.logger.Infow("running", "addr", service.conf.Portal.Gateway.Addr)
	return service.server.Run(ctx)
}
