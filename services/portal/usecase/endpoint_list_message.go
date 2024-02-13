package usecase

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type EndpointListMessageIn struct {
	*datastore.ScanningQuery
	WsId string
	EpId string
}

func (in *EndpointListMessageIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("ep_id", in.EpId, entities.IdNsEp),
	)
}

type EndpointListMessageOut struct {
	Data []entities.EndpointMessage
}

func (uc *endpoint) ListMessage(ctx context.Context, in *EndpointListMessageIn) (*EndpointListMessageOut, error) {
	ep, err := uc.repositories.Database().Endpoint().Get(ctx, in.WsId, in.EpId)
	if err != nil {
		return nil, err
	}

	reqMaps, err := uc.repositories.Datastore().Request().ScanMessages(ctx, ep.Id, in.ScanningQuery)
	if err != nil {
		return nil, err
	}

	resMaps, err := uc.repositories.Datastore().Response().ListMessages(ctx, ep.Id, reqMaps.MsgIds)
	if err != nil {
		return nil, err
	}

	msgses, err := uc.repositories.Datastore().Message().ListByIds(ctx, ep.AppId, reqMaps.MsgIds)
	if err != nil {
		return nil, err
	}

	out := &EndpointListMessageOut{}
	for _, msg := range msgses {
		out.Data = append(out.Data, *uc.mapping(reqMaps, resMaps, msg))
	}
	return out, nil
}
