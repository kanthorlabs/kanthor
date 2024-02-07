package usecase

import (
	"context"

	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"github.com/kanthorlabs/kanthor/services/sdk/repositories"
)

type Application interface {
	Create(ctx context.Context, in *ApplicationCreateIn) (*ApplicationCreateOut, error)
	Update(ctx context.Context, in *ApplicationUpdateIn) (*ApplicationUpdateOut, error)
	Delete(ctx context.Context, in *ApplicationDeleteIn) (*ApplicationDeleteOut, error)

	List(ctx context.Context, in *ApplicationListIn) (*ApplicationListOut, error)
	Get(ctx context.Context, in *ApplicationGetIn) (*ApplicationGetOut, error)
}

type application struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
