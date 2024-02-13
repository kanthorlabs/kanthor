package db

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"github.com/kanthorlabs/kanthor/internal/routing"
)

type Application interface {
	Scan(ctx context.Context, query *datastore.ScanningQuery) chan *datastore.ScanningRecord[[]entities.Application]
	GetRoutes(ctx context.Context, ids []string) (map[string][]routing.Route, error)
}
