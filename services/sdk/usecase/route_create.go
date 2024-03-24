package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/conductor"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/internal/database/scopes"
)

var ErrRouteCreate = errors.New("SDK.ROUTE.CREATE.ERROR")

func (uc *route) Create(ctx context.Context, in *RouteCreateIn) (*RouteCreateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// ensure that the endpoint exists in the requesting workspace before proceeding
	var ep entities.Endpoint
	err := uc.orm.WithContext(ctx).Scopes(scopes.UseEp(in.WsId)).
		Where(fmt.Sprintf("%s.id = ?", entities.TableEp), in.EpId).
		First(&ep).Error
	if err != nil {
		uc.logger.Errorw(ErrEndpointCreate.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrEndpointCreate
	}

	doc := &entities.Route{
		EpId:                ep.Id,
		Name:                in.Name,
		Priority:            in.Priority,
		Exclusionary:        in.Exclusionary,
		ConditionSource:     in.ConditionSource,
		ConditionExpression: in.ConditionExpression,
	}
	doc.SetId()
	doc.SetAuditFacttor(uc.watch.Now())

	if err := uc.orm.WithContext(ctx).Create(doc).Error; err != nil {
		uc.logger.Errorw(ErrRouteCreate.Error(), "error", err.Error(), "in", utils.Stringify(in), "route", utils.Stringify(doc))
		return nil, ErrRouteCreate
	}

	out := &RouteCreateOut{doc}
	return out, nil
}

type RouteCreateIn struct {
	WsId                string
	EpId                string
	Name                string
	Priority            int32
	Exclusionary        bool
	ConditionSource     string
	ConditionExpression string
}

func (in *RouteCreateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ROUTE.CREATE.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.ROUTE.CREATE.IN.EP_ID", in.EpId, entities.IdNsEp),
		validator.StringRequired("SDK.ROUTE.CREATE.IN.NAME", in.Name),
		validator.NumberInRange("SDK.ROUTE.CREATE.IN.PRIORITY", in.Priority, 1, 128),
		validator.Custom("SDK.ROUTE.CREATE.IN.CONDITION_SOURCE", &conductor.ConditionSource{Source: in.ConditionSource}),
		validator.Custom("SDK.ROUTE.CREATE.IN.CONDITION_SOURCE", &conductor.ConditionExression{Expression: in.ConditionExpression}),
	)
}

type RouteCreateOut struct {
	*entities.Route
}
