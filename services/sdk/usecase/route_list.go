package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/internal/database/scopes"
)

var ErrRouteList = errors.New("SDK.ROUTE.LIST.ERROR")

func (uc *route) List(ctx context.Context, in *RouteListIn) (*RouteListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	model := &entities.Route{}
	base := uc.orm.WithContext(ctx).Model(model).Scopes(scopes.UseRt(in.WsId))
	if in.EpId != "" {
		base = base.Where(fmt.Sprintf("%s.ep_id = ?", entities.TableRt), in.EpId)
	}

	var count int64
	var docs []*entities.Route

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
	WsId  string
	EpId  string
	Query *database.PagingQuery
}

func (in *RouteListIn) Validate() error {
	err := validator.Validate(
		validator.StringStartsWith("SDK.ROUTE.LIST.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWithIfNotEmpty("SDK.ROUTE.LIST.IN.EP_ID", in.EpId, entities.IdNsEp),
		validator.PointerNotNil("SDK.ROUTE.LIST.IN.QUERY", in.Query),
	)
	if err != nil {
		return err
	}

	if err := in.Query.Validate(); err != nil {
		return err
	}
	return nil
}

type RouteListOut struct {
	Count int64
	Data  []*entities.Route
}
