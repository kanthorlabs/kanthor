package db

import (
	"context"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Endpoint interface {
	Create(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error)
	Update(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error)
	Delete(ctx context.Context, doc *entities.Endpoint) error

	List(ctx context.Context, wsId, appId string, query *database.PagingQuery) ([]entities.Endpoint, error)
	Count(ctx context.Context, wsId, appId string, query *database.PagingQuery) (int64, error)
	Get(ctx context.Context, wsId string, id string) (*entities.Endpoint, error)
}

type EndpointRule interface {
	Create(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error)
	Update(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error)
	Delete(ctx context.Context, doc *entities.EndpointRule) error

	List(ctx context.Context, wsId, appId, epId string, query *database.PagingQuery) ([]entities.EndpointRule, error)
	Count(ctx context.Context, wsId, appId, epId string, query *database.PagingQuery) (int64, error)
	Get(ctx context.Context, wsId string, id string) (*entities.EndpointRule, error)
}
