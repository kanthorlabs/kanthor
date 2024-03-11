package entrypoint

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"github.com/kanthorlabs/kanthor/services/sdk/entrypoint/api"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

func NewApi(
	conf *config.Config,
	logger logging.Logger,
	infra infrastructure.Infrastructure,
	uc usecase.Sdk,
) (patterns.Runnable, error) {
	return api.New(conf, logger, infra, uc)
}
