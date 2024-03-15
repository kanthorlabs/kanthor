package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrApplicationCreate = errors.New("SDK.APPLICATION.CREATE.ERROR")

func (uc *application) Create(ctx context.Context, in *ApplicationCreateIn) (*ApplicationCreateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Application{
		WsId: in.WsId,
		Name: in.Name,
	}
	doc.SetId()
	doc.SetAuditFacttor(uc.watch.Now())

	if err := uc.orm.WithContext(ctx).Create(doc).Error; err != nil {
		uc.logger.Errorw(ErrApplicationCreate.Error(), "error", err.Error(), "in", utils.Stringify(in), "application", utils.Stringify(doc))
		return nil, ErrApplicationCreate
	}

	out := &ApplicationCreateOut{doc}
	return out, nil
}

type ApplicationCreateIn struct {
	WsId string
	Name string
}

func (in *ApplicationCreateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.APPLICATION.CREATE.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringRequired("SDK.APPLICATION.CREATE.IN.NAME", in.Name),
	)
}

type ApplicationCreateOut struct {
	*entities.Application
}
