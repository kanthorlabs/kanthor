package repositories

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

type Request interface {
	Save(ctx context.Context, docs []*entities.Request) error
}
