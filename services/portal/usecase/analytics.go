package usecase

import (
	"context"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/repositories"
)

type Analytics interface {
	GetOverview(ctx context.Context, in *AnalyticsGetOverviewIn) (*AnalyticsGetOverviewOut, error)
}

type analytics struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
