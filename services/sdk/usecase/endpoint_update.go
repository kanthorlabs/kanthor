package usecase

import (
	"context"
	"net/http"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type EndpointUpdateIn struct {
	WsId   string
	Id     string
	Name   string
	Method string
}

func (in *EndpointUpdateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsEp),
		validator.StringRequired("name", in.Name),
		validator.StringOneOf("method", in.Method, []string{http.MethodPost, http.MethodPut}),
	)
}

type EndpointUpdateOut struct {
	Doc *entities.Endpoint
}

func (uc *endpoint) Update(ctx context.Context, in *EndpointUpdateIn) (*EndpointUpdateOut, error) {
	ep, err := uc.repositories.Database().Endpoint().Get(ctx, in.WsId, in.Id)
	if err != nil {
		return nil, err
	}

	ep.Name = in.Name
	ep.Method = in.Method
	ep.SetAT(uc.infra.Timer.Now())
	doc, err := uc.repositories.Database().Endpoint().Update(ctx, ep)
	if err != nil {
		return nil, err
	}

	return &EndpointUpdateOut{Doc: doc}, nil
}
