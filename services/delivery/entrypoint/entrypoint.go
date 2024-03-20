package entrypoint

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/delivery/config"
	"github.com/kanthorlabs/kanthor/services/delivery/entrypoint/dispatcher"
	"github.com/kanthorlabs/kanthor/services/delivery/entrypoint/scheduler"
	"github.com/kanthorlabs/kanthor/services/delivery/usecase"
)

func NewScheduler(
	conf *config.Config,
	logger logging.Logger,
	infra infrastructure.Infrastructure,
	uc usecase.Delivery,
) (patterns.Runnable, error) {
	return scheduler.New(conf, logger, infra, uc)
}

func NewDispatcher(
	conf *config.Config,
	logger logging.Logger,
	infra infrastructure.Infrastructure,
	uc usecase.Delivery,
) (patterns.Runnable, error) {
	return dispatcher.New(conf, logger, infra, uc)
}
