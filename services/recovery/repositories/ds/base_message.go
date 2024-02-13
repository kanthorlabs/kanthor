package ds

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Message interface {
	Scan(ctx context.Context, appId string, query *datastore.ScanningQuery) chan *datastore.ScanningRecord[[]entities.Message]
}
