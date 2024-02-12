package usecase

import (
	"context"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/infrastructure/streaming"
	"github.com/kanthorlabs/kanthor/services/scheduler/config"
	"github.com/kanthorlabs/kanthor/services/scheduler/repositories"
)

type Request interface {
	Schedule(ctx context.Context, in *RequestScheduleIn) (*RequestScheduleOut, error)
}

type request struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	publisher    streaming.Publisher
	repositories repositories.Repositories
}
