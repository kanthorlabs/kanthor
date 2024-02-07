package ds

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Response interface {
	Create(ctx context.Context, docs []*entities.Response) ([]string, error)
}
