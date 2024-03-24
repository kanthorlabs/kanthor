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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrRouteUpdate = errors.New("SDK.ROUTE.UPDATE.ERROR")

func (uc *route) Update(ctx context.Context, in *RouteUpdateIn) (*RouteUpdateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Route{}
	err := uc.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Locking{Strength: clause.LockingStrengthShare}).
			Scopes(scopes.UseRt(in.WsId)).
			Where(fmt.Sprintf("%s.id = ?", entities.TableRt), in.Id).
			First(doc).Error
		if err != nil {
			uc.logger.Errorw(ErrRouteUpdate.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrRouteUpdate
		}

		doc.Name = in.Name
		doc.Priority = in.Priority
		doc.Exclusionary = in.Exclusionary
		doc.ConditionSource = in.ConditionSource
		doc.ConditionExpression = in.ConditionExpression
		doc.SetAuditFacttor(uc.watch.Now())
		if err := tx.Save(doc).Error; err != nil {
			uc.logger.Errorw(ErrRouteUpdate.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrRouteUpdate
		}

		return nil
	})
	if err != nil {
		// the error is already logged in the transaction
		return nil, err
	}

	out := &RouteUpdateOut{doc}
	return out, nil
}

type RouteUpdateIn struct {
	WsId                string
	Id                  string
	Name                string
	Priority            int32
	Exclusionary        bool
	ConditionSource     string
	ConditionExpression string
}

func (in *RouteUpdateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ROUTE.UPDATE.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.ROUTE.UPDATE.IN.ID", in.Id, entities.IdNsRt),
		validator.StringRequired("SDK.ROUTE.UPDATE.IN.NAME", in.Name),
		validator.NumberInRange("SDK.ROUTE.UPDATE.IN.PRIORITY", in.Priority, 1, 128),
		validator.Custom("SDK.ROUTE.UPDATE.IN.CONDITION_SOURCE", &conductor.ConditionSource{Source: in.ConditionSource}),
		validator.Custom("SDK.ROUTE.UPDATE.IN.CONDITION_SOURCE", &conductor.ConditionExression{Expression: in.ConditionExpression}),
	)
}

type RouteUpdateOut struct {
	*entities.Route
}
