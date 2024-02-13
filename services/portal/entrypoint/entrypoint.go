package entrypoint

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/entrypoint/rest"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
)

func Rest(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Portal,
) patterns.Runnable {
	return rest.New(conf, logger, infra, db, ds, uc)
}
