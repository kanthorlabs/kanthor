package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrApplicationList = errors.New("PORTAL.APPLICATION.LIST.ERROR")

func (uc *application) List(ctx context.Context, in *ApplicationListIn) (*ApplicationListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	var count int64
	var docs []*entities.Application

	model := &entities.Application{}
	base := uc.orm.Model(model).Where("ws_id = ?", in.WsId)

	if err := in.Query.SqlxCount(base, model.ColPrimary(), model.ColSearch()).Count(&count).Error; err != nil {
		return nil, ErrApplicationList
	}
	if err := in.Query.Sqlx(base, model.ColPrimary(), model.ColSearch()).Find(&docs).Error; err != nil {
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
		validator.StringStartsWith("SDK.WORKSPACE.LIST.IN.WS_ID", in.WsId, entities.IdNsWs),
	)
	if err != nil {
		return err
	}

	if in.Query != nil {
		return in.Query.Validate()
	}
	// must clone it to get the new one
	in.Query = database.DefaultPagingQuery.Clone()
	return nil
}

type ApplicationListOut struct {
	Count int64
	Data  []*entities.Application
}
