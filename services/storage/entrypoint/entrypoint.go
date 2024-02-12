package entrypoint

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/datastore"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/storage/config"
	"github.com/kanthorlabs/kanthor/services/storage/entrypoint/consumer"
	"github.com/kanthorlabs/kanthor/services/storage/usecase"
)

func Consumer(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	ds datastore.Datastore,
	uc usecase.Storage,
) (patterns.Runnable, error) {
	return consumer.New(conf, logger, infra, ds, uc)
}
