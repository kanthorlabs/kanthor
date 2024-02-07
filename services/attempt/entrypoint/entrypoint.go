package entrypoint

import (
	"github.com/kanthorlabs/kanthor/database"
	"github.com/kanthorlabs/kanthor/datastore"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/attempt/config"
	"github.com/kanthorlabs/kanthor/services/attempt/entrypoint/consumer"
	"github.com/kanthorlabs/kanthor/services/attempt/entrypoint/cronjob"
	"github.com/kanthorlabs/kanthor/services/attempt/entrypoint/endeavor"
	"github.com/kanthorlabs/kanthor/services/attempt/entrypoint/selector"
	"github.com/kanthorlabs/kanthor/services/attempt/entrypoint/trigger"
	"github.com/kanthorlabs/kanthor/services/attempt/usecase"
)

func Cronjob(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	return cronjob.New(conf, logger, infra, db, ds, uc)
}

func Consumer(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	return consumer.New(conf, logger, infra, db, ds, uc)
}

func Trigger(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	return trigger.New(conf, logger, infra, db, ds, uc)
}

func Selector(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	return selector.New(conf, logger, infra, db, ds, uc)
}

func Endeavor(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) patterns.Runnable {
	return endeavor.New(conf, logger, infra, db, ds, uc)
}
