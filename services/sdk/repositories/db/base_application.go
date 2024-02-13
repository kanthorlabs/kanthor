package db

import (
	"context"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Application interface {
	Create(ctx context.Context, doc *entities.Application) (*entities.Application, error)
	Update(ctx context.Context, doc *entities.Application) (*entities.Application, error)
	Delete(ctx context.Context, doc *entities.Application) error

	List(ctx context.Context, wsId string, query *database.PagingQuery) ([]entities.Application, error)
	Count(ctx context.Context, wsId string, query *database.PagingQuery) (int64, error)
	Get(ctx context.Context, wsId, id string) (*entities.Application, error)
}
