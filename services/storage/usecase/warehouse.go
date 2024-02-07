package usecase

import (
	"context"

	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/services/storage/config"
	"github.com/kanthorlabs/kanthor/services/storage/repositories"
)

type Warehouse interface {
	Put(ctx context.Context, in *WarehousePutIn) (*WarehousePutOut, error)
}

type warehose struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
