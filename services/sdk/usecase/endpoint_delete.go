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

var ErrEndpointDelete = errors.New("SDK.ENDPOINT.DELETE.ERROR")

func (uc *endpoint) Delete(ctx context.Context, in *EndpointDeleteIn) (*EndpointDeleteOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	joinstm := fmt.Sprintf("JOIN %s ON %s.id = %s.app_id", entities.TableApp, entities.TableApp, entities.TableEp)
	wherestm := fmt.Sprintf("%s.id = ? AND %s.id = ?", entities.TableApp, entities.TableEp)

	doc := &entities.Endpoint{}
	err := uc.orm.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Locking{Strength: clause.LockingStrengthShare}).
			InnerJoins(joinstm).Where(wherestm, in.AppId, in.Id).
			First(doc).Error
		if err != nil {
			uc.logger.Errorw(ErrEndpointDelete.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrEndpointDelete
		}

		doc.SetAuditFacttor(uc.watch.Now(), in.Modifier)
		if err := tx.Save(doc).Error; err != nil {
			uc.logger.Errorw(
				ErrEndpointDelete.Error(),
				"error", err.Error(),
				"in", utils.Stringify(in),
				"application", utils.Stringify(doc),
			)
			return ErrEndpointDelete
		}

		if err := tx.Delete(doc).Error; err != nil {
			uc.logger.Errorw(
				ErrEndpointDelete.Error(),
				"error", err.Error(),
				"in", utils.Stringify(in),
				"application", utils.Stringify(doc),
			)
			return ErrEndpointDelete
		}

		return nil
	})
	if err != nil {
		// the error is already logged in the transaction
		return nil, err
	}

	out := &EndpointDeleteOut{doc}
	return out, nil
}

type EndpointDeleteIn struct {
	Modifier string
	AppId    string
	Id       string
}

func (in *EndpointDeleteIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("SDK.ENDPOINT.UPDATE.IN.MODIFIER", in.Modifier),
		validator.StringStartsWith("SDK.ENDPOINT.CREATE.IN.APP_ID", in.AppId, entities.IdNsApp),
		validator.StringStartsWith("SDK.ENDPOINT.GET.IN.ID", in.Id, entities.IdNsEp),
	)
}

type EndpointDeleteOut struct {
	*entities.Endpoint
}
