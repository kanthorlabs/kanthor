package entrypoint

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/entrypoint/api"
)

func NewApi(
	conf *config.Config,
	logger logging.Logger,
) (patterns.Runnable, error) {
	return api.New(conf, logger)
}
