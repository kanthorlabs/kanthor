package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"github.com/kanthorlabs/kanthor/pkg/identifier"
)

type ApplicationCreateIn struct {
	WsId string
	Name string
}

func (in *ApplicationCreateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringRequired("name", in.Name),
	)
}

type ApplicationCreateOut struct {
	Doc *entities.Application
}

func (uc *application) Create(ctx context.Context, in *ApplicationCreateIn) (*ApplicationCreateOut, error) {
	doc := &entities.Application{
		WsId: in.WsId,
		Name: in.Name,
	}
	doc.Id = identifier.New(entities.IdNsApp)
	doc.SetAT(uc.infra.Timer.Now())

	app, err := uc.repositories.Database().Application().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	return &ApplicationCreateOut{Doc: app}, nil
}
