package entrypoint

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"github.com/kanthorlabs/kanthor/services/sdk/entrypoint/rest"
	"github.com/kanthorlabs/kanthor/services/sdk/usecase"
)

func Rest(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	uc usecase.Sdk,
) patterns.Runnable {
	return rest.New(conf, logger, infra, db, uc)
}
