package db

import (
	"context"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Application interface {
	CreateBatch(ctx context.Context, docs []entities.Application) ([]string, error)
	Count(ctx context.Context, wsId string, query *database.PagingQuery) (int64, error)
	Get(ctx context.Context, wsId, id string) (*entities.Application, error)
}
