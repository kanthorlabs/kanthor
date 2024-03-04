//go:build wireinject
// +build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/entrypoint"
)

func Portal(provider configuration.Provider) (patterns.Runnable, error) {
	wire.Build(
		config.New,
		logging.New,
		entrypoint.NewApi,
	)
	return nil, nil
}
