package repositories

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

type Request interface {
	Save(ctx context.Context, docs []*entities.Request) error
	Get(ctx context.Context, pks []entities.RequestPk) ([]*entities.Request, error)
}
