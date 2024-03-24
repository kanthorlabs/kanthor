package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/internal/database/scopes"
)

var ErrRouteGet = errors.New("SDK.ROUTE.GET.ERROR")

func (uc *route) Get(ctx context.Context, in *RouteGetIn) (*RouteGetOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Route{}
	err := uc.orm.WithContext(ctx).
		Scopes(scopes.UseRt(in.WsId)).
		Where(fmt.Sprintf("%s.id = ?", entities.TableRt), in.Id).
		First(doc).Error
	if err != nil {
		uc.logger.Errorw(ErrRouteGet.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrRouteGet
	}

	out := &RouteGetOut{doc}
	return out, nil
}

type RouteGetIn struct {
	WsId string
	Id   string
}

func (in *RouteGetIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ROUTE.GET.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.ROUTE.GET.IN.ID", in.Id, entities.IdNsRt),
	)
}

type RouteGetOut struct {
	*entities.Route
}
