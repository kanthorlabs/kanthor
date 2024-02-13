package ds

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Attempt interface {
	Scan(ctx context.Context, query *datastore.ScanningQuery, next int64, count int) chan *datastore.ScanningRecord[[]entities.Attempt]
	ListRequests(ctx context.Context, attempts map[string]*entities.Attempt) (map[string]*entities.Request, error)
	Update(ctx context.Context, updates map[string]*entities.AttemptState) map[string]error
}
