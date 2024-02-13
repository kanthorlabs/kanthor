//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/infrastructure"
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
