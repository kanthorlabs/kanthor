//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/kanthorlabs/kanthor/configuration"
	"github.com/kanthorlabs/kanthor/database"
	"github.com/kanthorlabs/kanthor/datastore"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/recovery/config"
	"github.com/kanthorlabs/kanthor/services/recovery/entrypoint"
	"github.com/kanthorlabs/kanthor/services/recovery/repositories"
	"github.com/kanthorlabs/kanthor/services/recovery/usecase"
)

func RecoveryCronjob(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		datastore.New,
		database.New,
		repositories.New,
		usecase.New,
		entrypoint.Cronjob,
	)
	return nil, nil
}

func RecoveryConsumer(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		datastore.New,
		database.New,
		repositories.New,
		usecase.New,
		entrypoint.Consumer,
	)
	return nil, nil
}
