package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrApplicationDelete = errors.New("SDK.APPLICATION.DELETE.ERROR")

func (uc *application) Delete(ctx context.Context, in *ApplicationDeleteIn) (*ApplicationDeleteOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Application{}
	err := uc.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("ws_id = ? AND id = ?", in.WsId, in.Id).
			First(doc).Error
		if err != nil {
			uc.logger.Errorw(ErrApplicationDelete.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrApplicationDelete
		}

		doc.SetAuditFacttor(uc.watch.Now())
		if err := tx.Save(doc).Error; err != nil {
			uc.logger.Errorw(
				ErrApplicationDelete.Error(),
				"error", err.Error(),
				"in", utils.Stringify(in),
				"application", utils.Stringify(doc),
			)
			return ErrApplicationDelete
		}

		if err := tx.Delete(doc).Error; err != nil {
			uc.logger.Errorw(
				ErrApplicationDelete.Error(),
				"error", err.Error(),
				"in", utils.Stringify(in),
				"application", utils.Stringify(doc),
			)
			return ErrApplicationDelete
		}

		return nil
	})
	if err != nil {
		// the error is already logged in the transaction
		return nil, err
	}

	out := &ApplicationDeleteOut{doc}
	return out, nil
}

type ApplicationDeleteIn struct {
	WsId string
	Id   string
}

func (in *ApplicationDeleteIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.APPLICATION.GET.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.APPLICATION.GET.IN.ID", in.Id, entities.IdNsApp),
	)
}

type ApplicationDeleteOut struct {
	*entities.Application
}
