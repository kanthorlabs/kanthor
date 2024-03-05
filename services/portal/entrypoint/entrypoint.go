package entrypoint

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/entrypoint/api"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

func NewApi(
	conf *config.Config,
	logger logging.Logger,
	infra infrastructure.Infrastructure,
	uc usecase.Portal,
) (patterns.Runnable, error) {
	return api.New(conf, logger, infra, uc)
}
