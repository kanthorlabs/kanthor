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

var ErrEndpointDelete = errors.New("SDK.ENDPOINT.DELETE.ERROR")

func (uc *endpoint) Delete(ctx context.Context, in *EndpointDeleteIn) (*EndpointDeleteOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Endpoint{}
	err := uc.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Locking{Strength: clause.LockingStrengthShare}).
			Scopes(scopes.UseEp(in.WsId)).
			Where(fmt.Sprintf("%s.id = ?", entities.TableEp), in.Id).
			First(doc).Error
		if err != nil {
			uc.logger.Errorw(ErrEndpointDelete.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrEndpointDelete
		}

		doc.SetAuditFacttor(uc.watch.Now())
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
	WsId string
	Id   string
}

func (in *EndpointDeleteIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ENDPOINT.GET.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.ENDPOINT.GET.IN.ID", in.Id, entities.IdNsEp),
	)
}

type EndpointDeleteOut struct {
	*entities.Endpoint
}
