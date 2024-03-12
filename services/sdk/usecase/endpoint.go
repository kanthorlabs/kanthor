package usecase

import (
	"context"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"gorm.io/gorm"
)

type Endpoint interface {
	Create(ctx context.Context, in *EndpointCreateIn) (*EndpointCreateOut, error)
	Get(ctx context.Context, in *EndpointGetIn) (*EndpointGetOut, error)
	List(ctx context.Context, in *EndpointListIn) (*EndpointListOut, error)
	Update(ctx context.Context, in *EndpointUpdateIn) (*EndpointUpdateOut, error)
	Delete(ctx context.Context, in *EndpointDeleteIn) (*EndpointDeleteOut, error)
}

type endpoint struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	orm    *gorm.DB
}
