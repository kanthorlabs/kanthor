package entrypoint

import (
	"github.com/kanthorlabs/kanthor/database"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/scheduler/config"
	"github.com/kanthorlabs/kanthor/services/scheduler/entrypoint/consumer"
	"github.com/kanthorlabs/kanthor/services/scheduler/usecase"
)

func Consumer(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	uc usecase.Scheduler,
) patterns.Runnable {
	return consumer.New(conf, logger, infra, db, uc)
}
