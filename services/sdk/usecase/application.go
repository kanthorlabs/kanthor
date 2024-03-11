package usecase

import (
	"context"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"gorm.io/gorm"
)

type Application interface {
	Create(ctx context.Context, in *ApplicationCreateIn) (*ApplicationCreateOut, error)
	// Get(ctx context.Context, in *ApplicationGetIn) (*ApplicationGetOut, error)
	// List(ctx context.Context, in *ApplicationListIn) (*ApplicationListOut, error)
	// Update(ctx context.Context, in *ApplicationUpdateIn) (*ApplicationUpdateOut, error)
	// Delete(ctx context.Context, in *ApplicationDeleteIn) (*ApplicationDeleteOut, error)
}

type application struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	orm    *gorm.DB
}
