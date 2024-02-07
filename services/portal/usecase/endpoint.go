package usecase

import (
	"context"

	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/repositories"
)

type Endpoint interface {
	ListMessage(ctx context.Context, in *EndpointListMessageIn) (*EndpointListMessageOut, error)
	GetMessage(ctx context.Context, in *EndpointGetMessageIn) (*EndpointGetMessageOut, error)
}

type endpoint struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
