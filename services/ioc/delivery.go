//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/delivery/config"
	"github.com/kanthorlabs/kanthor/services/delivery/entrypoint"
	"github.com/kanthorlabs/kanthor/services/delivery/usecase"
)

func Scheduler(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		clock.New,
		wire.FieldsOf(new(*config.Config), "Infrastructure"),
		infrastructure.New,
		usecase.New,
		entrypoint.NewScheduler,
	)
	return nil, nil
}

func Dispatcher(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		clock.New,
		wire.FieldsOf(new(*config.Config), "Infrastructure"),
		infrastructure.New,
		usecase.New,
		entrypoint.NewDispatcher,
	)
	return nil, nil
}
