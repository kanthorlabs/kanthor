package usecase

import (
	"context"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"gorm.io/gorm"
)

type Route interface {
	Create(ctx context.Context, in *RouteCreateIn) (*RouteCreateOut, error)
	Get(ctx context.Context, in *RouteGetIn) (*RouteGetOut, error)
	List(ctx context.Context, in *RouteListIn) (*RouteListOut, error)
	Update(ctx context.Context, in *RouteUpdateIn) (*RouteUpdateOut, error)
	Delete(ctx context.Context, in *RouteDeleteIn) (*RouteDeleteOut, error)
}

type route struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	orm    *gorm.DB
}
