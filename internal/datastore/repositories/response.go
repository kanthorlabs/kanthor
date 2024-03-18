package repositories

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

type Response interface {
	Save(ctx context.Context, docs []*entities.Response) error
}
