//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/scheduler/config"
	"github.com/kanthorlabs/kanthor/services/scheduler/entrypoint"
	"github.com/kanthorlabs/kanthor/services/scheduler/repositories"
	"github.com/kanthorlabs/kanthor/services/scheduler/usecase"
)

func Scheduler(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		database.New,
		repositories.New,
		usecase.New,
		entrypoint.Consumer,
	)
	return nil, nil
}
