package entrypoint

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/storage/config"
	"github.com/kanthorlabs/kanthor/services/storage/entrypoint/consumer"
	"github.com/kanthorlabs/kanthor/services/storage/usecase"
)

func NewConsumer(
	conf *config.Config,
	logger logging.Logger,
	infra infrastructure.Infrastructure,
	uc usecase.Storage,
) (patterns.Runnable, error) {
	return consumer.New(conf, logger, infra, uc)
}
