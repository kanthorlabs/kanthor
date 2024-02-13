package usecase

import (
	"context"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type ApplicationListIn struct {
	*database.PagingQuery
	WsId string
}

func (in *ApplicationListIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
	)
}

type ApplicationListOut struct {
	Data  []entities.Application
	Count int64
}

func (uc *application) List(ctx context.Context, in *ApplicationListIn) (*ApplicationListOut, error) {
	data, err := uc.repositories.Database().Application().List(ctx, in.WsId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	count, err := uc.repositories.Database().Application().Count(ctx, in.WsId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	out := &ApplicationListOut{Data: data, Count: count}
	return out, nil
}
