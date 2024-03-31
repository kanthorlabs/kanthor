package repositories

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

type Response interface {
	Save(ctx context.Context, docs []*entities.Response) error
	Get(ctx context.Context, pks []entities.ResponsePk) ([]*entities.Response, error)
}
