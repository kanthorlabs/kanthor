package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrRouteDelete = errors.New("SDK.ROUTE.DELETE.ERROR")

func (uc *route) Delete(ctx context.Context, in *RouteDeleteIn) (*RouteDeleteOut, error) {
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
			uc.logger.Errorw(ErrRouteDelete.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrRouteDelete
		}

		doc.SetAuditFacttor(in.Modifier, uc.watch.Now())
		if err := tx.Save(doc).Error; err != nil {
			uc.logger.Errorw(
				ErrRouteDelete.Error(),
				"error", err.Error(),
				"in", utils.Stringify(in),
				"application", utils.Stringify(doc),
			)
			return ErrRouteDelete
		}

		if err := tx.Delete(doc).Error; err != nil {
			uc.logger.Errorw(
				ErrRouteDelete.Error(),
				"error", err.Error(),
				"in", utils.Stringify(in),
				"application", utils.Stringify(doc),
			)
			return ErrRouteDelete
		}

		return nil
	})
	if err != nil {
		// the error is already logged in the transaction
		return nil, err
	}

	out := &RouteDeleteOut{doc}
	return out, nil
}

type RouteDeleteIn struct {
	Modifier string
	EpId     string
	Id       string
}

func (in *RouteDeleteIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("SDK.ROUTE.DELETE.IN.MODIFIER", in.Modifier),
		validator.StringStartsWith("SDK.ROUTE.DELETE.IN.EP_ID", in.EpId, entities.IdNsEp),
		validator.StringStartsWith("SDK.ROUTE.DELETE.IN.ID", in.Id, entities.IdNsRt),
	)
}

type RouteDeleteOut struct {
	*entities.Route
}
