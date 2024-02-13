package ds

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Request interface {
	Scan(ctx context.Context, epId string, query *datastore.ScanningQuery) chan *datastore.ScanningRecord[[]entities.Request]
}
