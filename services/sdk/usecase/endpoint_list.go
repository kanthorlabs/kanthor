package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrEndpointList = errors.New("SDK.ENDPOINT.LIST.ERROR")

func (uc *endpoint) List(ctx context.Context, in *EndpointListIn) (*EndpointListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	var count int64
	var docs []*entities.Endpoint

	model := &entities.Endpoint{}
	base := uc.orm.WithContext(ctx).Model(model).Where("app_id = ?", in.AppId)

	if err := in.Query.SqlxCount(base, model.PrimaryProp(), model.SearchProps()).Count(&count).Error; err != nil {
		uc.logger.Errorw(ErrEndpointList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrEndpointList
	}
	if err := in.Query.Sqlx(base, model.PrimaryProp(), model.SearchProps()).Find(&docs).Error; err != nil {
		uc.logger.Errorw(ErrEndpointList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrEndpointList
	}

	return &EndpointListOut{Count: count, Data: docs}, nil
}

type EndpointListIn struct {
	AppId string
	Query *database.PagingQuery
}

func (in *EndpointListIn) Validate() error {
	err := validator.Validate(
		validator.StringStartsWith("SDK.ENDPOINT.LIST.IN.APP_ID", in.AppId, entities.IdNsApp),
		validator.PointerNotNil("SDK.ENDPOINT.LIST.IN.QUERY", in.Query),
	)
	if err != nil {
		return err
	}

	if err := in.Query.Validate(); err != nil {
		return err
	}
	return nil
}

type EndpointListOut struct {
	Count int64
	Data  []*entities.Endpoint
}
