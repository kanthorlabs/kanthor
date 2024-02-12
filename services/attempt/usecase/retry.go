package usecase

import (
	"context"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/infrastructure/streaming"
	"github.com/kanthorlabs/kanthor/services/attempt/config"
	"github.com/kanthorlabs/kanthor/services/attempt/repositories"
)

type Retry interface {
	Trigger(ctx context.Context, in *RetryTriggerIn) (*RetryTriggerOut, error)
	Select(ctx context.Context, in *RetrySelectIn) (*RetrySelectOut, error)
	Endeavor(ctx context.Context, in *RetryEndeavorIn) (*RetryEndeavorOut, error)
}

type retry struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	publisher    streaming.Publisher
	repositories repositories.Repositories
}
