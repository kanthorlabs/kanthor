package usecase

import (
	"context"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/infrastructure/streaming"
	"github.com/kanthorlabs/kanthor/services/recovery/config"
	"github.com/kanthorlabs/kanthor/services/recovery/repositories"
)

type Scanner interface {
	Schedule(ctx context.Context, in *ScannerScheduleIn) (*ScannerScheduleOut, error)
	Execute(ctx context.Context, in *ScannerExecuteIn) (*ScannerExecuteOut, error)
}

type scanner struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	publisher    streaming.Publisher
	repositories repositories.Repositories
}
