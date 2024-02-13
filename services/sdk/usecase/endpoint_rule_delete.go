package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type EndpointRuleDeleteIn struct {
	WsId string
	Id   string
}

func (in *EndpointRuleDeleteIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsEpr),
	)
}

type EndpointRuleDeleteOut struct {
	Doc *entities.EndpointRule
}

func (uc *endpointRule) Delete(ctx context.Context, in *EndpointRuleDeleteIn) (*EndpointRuleDeleteOut, error) {
	epr, err := uc.repositories.Database().EndpointRule().Get(ctx, in.WsId, in.Id)
	if err != nil {
		return nil, err
	}

	if err := uc.repositories.Database().EndpointRule().Delete(ctx, epr); err != nil {
		return nil, err
	}

	return &EndpointRuleDeleteOut{Doc: epr}, nil
}
