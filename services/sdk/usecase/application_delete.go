package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type ApplicationDeleteIn struct {
	WsId string
	Id   string
}

func (in *ApplicationDeleteIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", in.Id, entities.IdNsApp),
	)
}

type ApplicationDeleteOut struct {
	Doc *entities.Application
}

func (uc *application) Delete(ctx context.Context, in *ApplicationDeleteIn) (*ApplicationDeleteOut, error) {
	app, err := uc.repositories.Database().Application().Get(ctx, in.WsId, in.Id)
	if err != nil {
		return nil, err
	}

	if err := uc.repositories.Database().Application().Delete(ctx, app); err != nil {
		return nil, err
	}

	return &ApplicationDeleteOut{Doc: app}, nil
}
