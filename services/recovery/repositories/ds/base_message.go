package ds

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Message interface {
	Scan(ctx context.Context, appId string, query *entities.ScanningQuery) chan *entities.ScanningResult[[]entities.Message]
}
