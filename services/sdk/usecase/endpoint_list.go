package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type EndpointListIn struct {
	*entities.PagingQuery
	WsId  string
	AppId string
}

func (in *EndpointListIn) Validate() error {
	if err := in.PagingQuery.Validate(); err != nil {
		return err
	}

	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWithIfNotEmpty("app_id", in.AppId, entities.IdNsApp),
	)
}

type EndpointListOut struct {
	Data  []entities.Endpoint
	Count int64
}

func (uc *endpoint) List(ctx context.Context, in *EndpointListIn) (*EndpointListOut, error) {
	data, err := uc.repositories.Database().Endpoint().List(ctx, in.WsId, in.AppId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	count, err := uc.repositories.Database().Endpoint().Count(ctx, in.WsId, in.AppId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	out := &EndpointListOut{Data: data, Count: count}
	return out, nil
}
