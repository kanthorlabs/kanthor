package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type ApplicationGetMessageIn struct {
	WsId  string
	AppId string
	Id    string
}

func (in *ApplicationGetMessageIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
		validator.StringStartsWith("id", in.Id, entities.IdNsMsg),
	)
}

type ApplicationGetMessageOut struct {
	Doc *entities.Message
}

func (uc *application) GetMessage(ctx context.Context, in *ApplicationGetMessageIn) (*ApplicationGetMessageOut, error) {
	app, err := uc.repositories.Database().Application().Get(ctx, in.WsId, in.AppId)
	if err != nil {
		return nil, err
	}

	msg, err := uc.repositories.Datastore().Message().Get(ctx, app.Id, in.Id)
	if err != nil {
		return nil, err
	}

	out := &ApplicationGetMessageOut{Doc: msg}
	return out, nil
}
