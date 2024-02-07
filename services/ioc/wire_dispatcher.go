//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/kanthorlabs/kanthor/configuration"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/dispatcher/config"
	"github.com/kanthorlabs/kanthor/services/dispatcher/entrypoint"
	"github.com/kanthorlabs/kanthor/services/dispatcher/usecase"
)

func Dispatcher(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		usecase.New,
		entrypoint.Consumer,
	)
	return nil, nil
}
