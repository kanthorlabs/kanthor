package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type ApplicationUpdateIn struct {
	WsId string
	Id   string
	Name string
}

func (in *ApplicationUpdateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsApp),
		validator.StringRequired("name", in.Name),
	)
}

type ApplicationUpdateOut struct {
	Doc *entities.Application
}

func (uc *application) Update(ctx context.Context, in *ApplicationUpdateIn) (*ApplicationUpdateOut, error) {
	app, err := uc.repositories.Database().Application().Get(ctx, in.WsId, in.Id)
	if err != nil {
		return nil, err
	}

	app.Name = in.Name
	app.SetAT(uc.infra.Timer.Now())
	doc, err := uc.repositories.Database().Application().Update(ctx, app)
	if err != nil {
		return nil, err
	}

	return &ApplicationUpdateOut{Doc: doc}, nil
}
