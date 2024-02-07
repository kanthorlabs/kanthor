package usecase

import (
	"context"

	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/repositories"
)

type Application interface {
	ListMessage(ctx context.Context, in *ApplicationListMessageIn) (*ApplicationListMessageOut, error)
	GetMessage(ctx context.Context, in *ApplicationGetMessageIn) (*ApplicationGetMessageOut, error)
}

type application struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
