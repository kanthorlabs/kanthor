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
	"github.com/kanthorlabs/kanthor/services/attempt/config"
	"github.com/kanthorlabs/kanthor/services/attempt/entrypoint"
	"github.com/kanthorlabs/kanthor/services/attempt/repositories"
	"github.com/kanthorlabs/kanthor/services/attempt/usecase"
)

func AttemptCronjob(provider configuration.Provider) (patterns.Runnable, error) {
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

func AttemptConsumer(provider configuration.Provider) (patterns.Runnable, error) {
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

func AttemptTrigger(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		datastore.New,
		database.New,
		repositories.New,
		usecase.New,
		entrypoint.Trigger,
	)
	return nil, nil
}

func AttemptSelector(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		datastore.New,
		database.New,
		repositories.New,
		usecase.New,
		entrypoint.Selector,
	)
	return nil, nil
}

func AttemptEndeavor(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		datastore.New,
		database.New,
		repositories.New,
		usecase.New,
		entrypoint.Endeavor,
	)
	return nil, nil
}
