package ds

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Attempt interface {
	Create(ctx context.Context, docs []*entities.Attempt) ([]string, error)
}
