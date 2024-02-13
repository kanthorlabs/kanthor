package entrypoint

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/infrastructure"
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
) (patterns.Runnable, error) {
	return cronjob.New(conf, logger, infra, db, ds, uc)
}

func Consumer(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) (patterns.Runnable, error) {
	return consumer.New(conf, logger, infra, db, ds, uc)
}

func Trigger(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) (patterns.Runnable, error) {
	return trigger.New(conf, logger, infra, db, ds, uc)
}

func Selector(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) (patterns.Runnable, error) {
	return selector.New(conf, logger, infra, db, ds, uc)
}

func Endeavor(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Attempt,
) (patterns.Runnable, error) {
	return endeavor.New(conf, logger, infra, db, ds, uc)
}
