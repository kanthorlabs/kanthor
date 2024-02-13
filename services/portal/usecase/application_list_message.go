package usecase

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type ApplicationListMessageIn struct {
	*datastore.ScanningQuery
	WsId  string
	AppId string
}

func (in *ApplicationListMessageIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
	)
}

type ApplicationListMessageOut struct {
	Data []entities.Message
}

func (uc *application) ListMessage(ctx context.Context, in *ApplicationListMessageIn) (*ApplicationListMessageOut, error) {
	app, err := uc.repositories.Database().Application().Get(ctx, in.WsId, in.AppId)
	if err != nil {
		return nil, err
	}

	data, err := uc.repositories.Datastore().Message().Scan(ctx, app.Id, in.ScanningQuery)
	if err != nil {
		return nil, err
	}

	out := &ApplicationListMessageOut{Data: data}
	return out, nil
}
