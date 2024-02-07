package entrypoint

import (
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/dispatcher/config"
	"github.com/kanthorlabs/kanthor/services/dispatcher/entrypoint/consumer"
	"github.com/kanthorlabs/kanthor/services/dispatcher/usecase"
)

func Consumer(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	uc usecase.Dispatcher,
) patterns.Runnable {
	return consumer.New(conf, logger, infra, uc)
}
