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

var ErrApplicationUpdate = errors.New("PORTAL.APPLICATION.UPDATE.ERROR")

func (uc *application) Update(ctx context.Context, in *ApplicationUpdateIn) (*ApplicationUpdateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Application{}

	err := uc.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Locking{Strength: clause.LockingStrengthShare}).
			Where("ws_id = ? AND id = ?", in.WsId, in.Id).
			First(doc).Error
		if err != nil {
			uc.logger.Errorw(ErrApplicationUpdate.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrApplicationUpdate
		}

		doc.Name = in.Name
		doc.SetAuditFacttor(uc.watch.Now())
		if err := tx.Save(doc).Error; err != nil {
			uc.logger.Errorw(ErrApplicationUpdate.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrApplicationUpdate
		}

		return nil
	})
	if err != nil {
		// the error is already logged in the transaction
		return nil, err
	}

	out := &ApplicationUpdateOut{doc}
	return out, nil
}

type ApplicationUpdateIn struct {
	WsId string
	Id   string
	Name string
}

func (in *ApplicationUpdateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.APPLICATION.GET.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.APPLICATION.GET.IN.ID", in.Id, entities.IdNsApp),
		validator.StringRequired("SDK.APPLICATION.UPDATE.IN.NAME", in.Name),
	)
}

type ApplicationUpdateOut struct {
	*entities.Application
}
