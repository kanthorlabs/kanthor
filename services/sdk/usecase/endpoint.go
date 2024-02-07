package usecase

import (
	"context"

	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"github.com/kanthorlabs/kanthor/services/sdk/repositories"
)

type Endpoint interface {
	Create(ctx context.Context, in *EndpointCreateIn) (*EndpointCreateOut, error)
	Update(ctx context.Context, in *EndpointUpdateIn) (*EndpointUpdateOut, error)
	Delete(ctx context.Context, in *EndpointDeleteIn) (*EndpointDeleteOut, error)

	List(ctx context.Context, in *EndpointListIn) (*EndpointListOut, error)
	Get(ctx context.Context, in *EndpointGetIn) (*EndpointGetOut, error)
}

type endpoint struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
