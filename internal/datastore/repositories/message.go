package repositories

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

type Message interface {
	Save(ctx context.Context, docs []*entities.Message) error
	Get(ctx context.Context, pks []entities.MessagePk) ([]*entities.Message, error)
}
