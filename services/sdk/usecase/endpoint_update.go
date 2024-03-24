package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/internal/database/scopes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrEndpointUpdate = errors.New("SDK.ENDPOINT.UPDATE.ERROR")

func (uc *endpoint) Update(ctx context.Context, in *EndpointUpdateIn) (*EndpointUpdateOut, error) {
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
			uc.logger.Errorw(ErrEndpointUpdate.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrEndpointUpdate
		}

		doc.Name = in.Name
		doc.Method = in.Method
		doc.Uri = in.Uri
		doc.SetAuditFacttor(uc.watch.Now())
		if err := tx.Save(doc).Error; err != nil {
			uc.logger.Errorw(ErrEndpointUpdate.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrEndpointUpdate
		}

		return nil
	})
	if err != nil {
		// the error is already logged in the transaction
		return nil, err
	}

	out := &EndpointUpdateOut{doc}
	return out, nil
}

type EndpointUpdateIn struct {
	WsId   string
	Id     string
	Name   string
	Method string
	Uri    string
}

func (in *EndpointUpdateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ENDPOINT.UPDATE.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.ENDPOINT.UPDATE.IN.ID", in.Id, entities.IdNsEp),
		validator.StringRequired("SDK.ENDPOINT.UPDATE.IN.NAME", in.Name),
		validator.StringOneOf("SDK.ENDPOINT.UPDATE.IN.METHOD", in.Method, []string{http.MethodPost, http.MethodPut}),
		validator.StringUri("SDK.ENDPOINT.UPDATE.IN.URI", in.Uri),
	)
}

type EndpointUpdateOut struct {
	*entities.Endpoint
}
