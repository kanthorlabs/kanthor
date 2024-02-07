package ds

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Request interface {
	Create(ctx context.Context, docs []*entities.Request) ([]string, error)
}
