package repositories

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

type Message interface {
	Save(ctx context.Context, docs []*entities.Message) error
}
