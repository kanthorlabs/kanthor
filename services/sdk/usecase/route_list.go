package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrRouteList = errors.New("SDK.ROUTE.LIST.ERROR")

func (uc *route) List(ctx context.Context, in *RouteListIn) (*RouteListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	var count int64
	var docs []*entities.Route

	model := &entities.Route{}
	base := uc.orm.Model(model).Where("ep_id = ?", in.EpId)

	if err := in.Query.SqlxCount(base, model.PrimaryProp(), model.SearchProps()).Count(&count).Error; err != nil {
		uc.logger.Errorw(ErrRouteList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrRouteList
	}
	if err := in.Query.Sqlx(base, model.PrimaryProp(), model.SearchProps()).Find(&docs).Error; err != nil {
		uc.logger.Errorw(ErrRouteList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrRouteList
	}

	return &RouteListOut{Count: count, Data: docs}, nil
}

type RouteListIn struct {
	EpId  string
	Query *database.PagingQuery
}

func (in *RouteListIn) Validate() error {
	err := validator.Validate(
		validator.StringStartsWith("SDK.ROUTE.LIST.IN.EP_ID", in.EpId, entities.IdNsEp),
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

type RouteListOut struct {
	Count int64
	Data  []*entities.Route
}
