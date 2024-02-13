package entrypoint

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/recovery/config"
	"github.com/kanthorlabs/kanthor/services/recovery/entrypoint/consumer"
	"github.com/kanthorlabs/kanthor/services/recovery/entrypoint/cronjob"
	"github.com/kanthorlabs/kanthor/services/recovery/usecase"
)

func Cronjob(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Recovery,
) (patterns.Runnable, error) {
	return cronjob.New(conf, logger, infra, db, ds, uc)
}

func Consumer(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Recovery,
) (patterns.Runnable, error) {
	return consumer.New(conf, logger, infra, db, ds, uc)
}
