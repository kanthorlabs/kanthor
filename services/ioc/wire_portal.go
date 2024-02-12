//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/database"
	"github.com/kanthorlabs/kanthor/datastore"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/entrypoint"
	"github.com/kanthorlabs/kanthor/services/portal/repositories"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

func Portal(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		infrastructure.New,
		database.New,
		datastore.New,
		repositories.New,
		usecase.New,
		entrypoint.Rest,
	)
	return nil, nil
}
