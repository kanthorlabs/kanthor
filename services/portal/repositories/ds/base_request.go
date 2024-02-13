package ds

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type MessageRequestMaps struct {
	Maps   map[string][]entities.Request
	MsgIds []string
}

type Request interface {
	ScanMessages(ctx context.Context, epId string, query *datastore.ScanningQuery) (*MessageRequestMaps, error)
	GetMessage(ctx context.Context, epId, msgId string) (*MessageRequestMaps, error)
	Scan(ctx context.Context, epId string, query *datastore.ScanningQuery) ([]entities.Request, error)
}
