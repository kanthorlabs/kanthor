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
	"github.com/kanthorlabs/kanthor/services/storage/config"
	"github.com/kanthorlabs/kanthor/services/storage/entrypoint"
	"github.com/kanthorlabs/kanthor/services/storage/usecase"
)

func Storage(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		clock.New,
		wire.FieldsOf(new(*config.Config), "Infrastructure"),
		infrastructure.New,
		usecase.New,
		entrypoint.NewConsumer,
	)
	return nil, nil
}
