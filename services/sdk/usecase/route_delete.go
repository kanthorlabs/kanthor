package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/internal/database/scopes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrRouteDelete = errors.New("SDK.ROUTE.DELETE.ERROR")

func (uc *route) Delete(ctx context.Context, in *RouteDeleteIn) (*RouteDeleteOut, error) {
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
			uc.logger.Errorw(ErrRouteDelete.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrRouteDelete
		}

		doc.SetAuditFacttor(uc.watch.Now())
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
	WsId string
	Id   string
}

func (in *RouteDeleteIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ROUTE.UPDATE.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.ROUTE.DELETE.IN.ID", in.Id, entities.IdNsRt),
	)
}

type RouteDeleteOut struct {
	*entities.Route
}
