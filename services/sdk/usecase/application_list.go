package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrApplicationList = errors.New("SDK.APPLICATION.LIST.ERROR")

func (uc *application) List(ctx context.Context, in *ApplicationListIn) (*ApplicationListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	model := &entities.Application{}
	base := uc.orm.WithContext(ctx).Model(model).Where("ws_id = ?", in.WsId)

	var count int64
	var docs []*entities.Application

	if err := in.Query.SqlxCount(base, model.PrimaryProp(), model.SearchProps()).Count(&count).Error; err != nil {
		uc.logger.Errorw(ErrApplicationList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrApplicationList
	}
	if err := in.Query.Sqlx(base, model.PrimaryProp(), model.SearchProps()).Find(&docs).Error; err != nil {
		uc.logger.Errorw(ErrApplicationList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrApplicationList
	}

	return &ApplicationListOut{Count: count, Data: docs}, nil
}

type ApplicationListIn struct {
	WsId  string
	Query *database.PagingQuery
}

func (in *ApplicationListIn) Validate() error {
	err := validator.Validate(
		validator.StringStartsWith("SDK.APPLICATION.LIST.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.PointerNotNil("SDK.APPLICATION.LIST.IN.QUERY", in.Query),
	)
	if err != nil {
		return err
	}

	if err := in.Query.Validate(); err != nil {
		return err
	}
	return nil
}

type ApplicationListOut struct {
	Count int64
	Data  []*entities.Application
}
