package ds

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Message interface {
	Create(ctx context.Context, docs []*entities.Message) ([]string, error)
}
