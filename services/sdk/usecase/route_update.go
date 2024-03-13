package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/conductor"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrRouteUpdate = errors.New("SDK.ROUTE.UPDATE.ERROR")

func (uc *route) Update(ctx context.Context, in *RouteUpdateIn) (*RouteUpdateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	joinstm := fmt.Sprintf("JOIN %s ON %s.id = %s.ep_id", entities.TableEp, entities.TableEp, entities.TableRt)
	wherestm := fmt.Sprintf("%s.id = ? AND %s.id = ?", entities.TableEp, entities.TableRt)

	doc := &entities.Route{}
	err := uc.orm.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Locking{Strength: clause.LockingStrengthShare}).
			InnerJoins(joinstm).Where(wherestm, in.EpId, in.Id).
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
		doc.SetAuditFacttor(in.Modifier, uc.watch.Now())
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
	Modifier            string
	EpId                string
	Id                  string
	Name                string
	Priority            int32
	Exclusionary        bool
	ConditionSource     string
	ConditionExpression string
}

func (in *RouteUpdateIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("SDK.ROUTE.UPDATE.IN.MODIFIER", in.Modifier),
		validator.StringStartsWith("SDK.ROUTE.UPDATE.IN.EP_ID", in.EpId, entities.IdNsEp),
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
