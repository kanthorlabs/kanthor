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
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"github.com/kanthorlabs/kanthor/services/sdk/entrypoint"
	"github.com/kanthorlabs/kanthor/services/sdk/repositories"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

func Sdk(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		database.New,
		repositories.New,
		usecase.New,
		entrypoint.Rest,
	)
	return nil, nil
}
