package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrApplicationGet = errors.New("PORTAL.APPLICATION.GET.ERROR")

func (uc *application) Get(ctx context.Context, in *ApplicationGetIn) (*ApplicationGetOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Application{}
	err := uc.orm.Where("ws_id = ? AND id = ?", in.WsId, in.Id).First(doc).Error
	if err != nil {
		uc.logger.Errorw(ErrApplicationGet.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrApplicationGet
	}

	out := &ApplicationGetOut{doc}
	return out, nil
}

type ApplicationGetIn struct {
	WsId string
	Id   string
}

func (in *ApplicationGetIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.APPLICATION.GET.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.APPLICATION.GET.IN.ID", in.Id, entities.IdNsApp),
	)
}

type ApplicationGetOut struct {
	*entities.Application
}
