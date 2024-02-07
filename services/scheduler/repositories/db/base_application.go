package db

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/entities"
	"github.com/kanthorlabs/kanthor/internal/routing"
)

type Application interface {
	Get(ctx context.Context, id string) (*entities.Application, error)
	GetRoutes(ctx context.Context, ids []string) (map[string][]routing.Route, error)
}
