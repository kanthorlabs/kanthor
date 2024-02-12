package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type EndpointGetIn struct {
	WsId string
	Id   string
}

func (in *EndpointGetIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsEp),
	)
}

type EndpointGetOut struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Get(ctx context.Context, in *EndpointGetIn) (*EndpointGetOut, error) {
	ep, err := uc.repositories.Database().Endpoint().Get(ctx, in.WsId, in.Id)
	if err != nil {
		return nil, err
	}
	return &EndpointGetOut{Doc: ep}, nil
}
