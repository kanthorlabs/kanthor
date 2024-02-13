package db

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Endpoint interface {
	Scan(ctx context.Context, query *datastore.ScanningQuery) chan *datastore.ScanningRecord[[]entities.Endpoint]
}
