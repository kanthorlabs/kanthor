package usecase

import (
	"context"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type EndpointRuleListIn struct {
	*database.PagingQuery
	WsId  string
	AppId string
	EpId  string
}

func (in *EndpointRuleListIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWithIfNotEmpty("app_id", in.AppId, entities.IdNsApp),
		validator.StringStartsWithIfNotEmpty("ep_id", in.EpId, entities.IdNsEp),
	)
}

type EndpointRuleListOut struct {
	Data  []entities.EndpointRule
	Count int64
}

func (uc *endpointRule) List(ctx context.Context, in *EndpointRuleListIn) (*EndpointRuleListOut, error) {
	data, err := uc.repositories.Database().EndpointRule().List(ctx, in.WsId, in.AppId, in.EpId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	count, err := uc.repositories.Database().EndpointRule().Count(ctx, in.WsId, in.AppId, in.EpId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	out := &EndpointRuleListOut{Data: data, Count: count}
	return out, nil
}
