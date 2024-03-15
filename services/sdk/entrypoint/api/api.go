package api

import (
	"context"
	"errors"
	"sync"

	"github.com/kanthorlabs/common/gateway"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra infrastructure.Infrastructure,
	uc usecase.Sdk,
) (patterns.Runnable, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	entrypoint := &sdk{
		conf:   conf,
		logger: logger.With("entrypoint", "api"),
		infra:  infra,
		uc:     uc,
	}
	return entrypoint, nil
}

type sdk struct {
	conf   *config.Config
	logger logging.Logger
	infra  infrastructure.Infrastructure
	uc     usecase.Sdk

	server gateway.Gateway
	mu     sync.Mutex
	status int
}

func (service *sdk) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status == patterns.StatusStarted {
		return ErrAlreadyStarted
	}

	if err := service.infra.Connect(ctx); err != nil {
		return err
	}

	server, err := gateway.New(&service.conf.Sdk.Gateway, service.logger)
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

func (service *sdk) Stop(ctx context.Context) error {
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

	if err := service.infra.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	service.logger.Info("stopped")
	return returning
}

func (service *sdk) Run(ctx context.Context) error {
	if err := service.httpx(); err != nil {
		return err
	}

	service.logger.Infow("running", "addr", service.conf.Sdk.Gateway.Addr)
	return service.server.Run(ctx)
}
