package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrEndpointUpdate = errors.New("SDK.ENDPOINT.UPDATE.ERROR")

func (uc *endpoint) Update(ctx context.Context, in *EndpointUpdateIn) (*EndpointUpdateOut, error) {
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
			uc.logger.Errorw(ErrEndpointUpdate.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrEndpointUpdate
		}

		doc.Name = in.Name
		doc.Method = in.Method
		doc.Uri = in.Uri
		doc.SetAuditFacttor(uc.watch.Now(), in.Modifier)
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
	Modifier string
	AppId    string
	Id       string
	Name     string
	Method   string
	Uri      string
}

func (in *EndpointUpdateIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("SDK.ENDPOINT.UPDATE.IN.MODIFIER", in.Modifier),
		validator.StringStartsWith("SDK.ENDPOINT.CREATE.IN.APP_ID", in.AppId, entities.IdNsApp),
		validator.StringStartsWith("SDK.ENDPOINT.GET.IN.ID", in.Id, entities.IdNsEp),
		validator.StringRequired("SDK.ENDPOINT.UPDATE.IN.NAME", in.Name),
		validator.StringOneOf("SDK.ENDPOINT.CREATE.IN.METHOD", in.Method, []string{http.MethodPost, http.MethodPut}),
		validator.StringUri("SDK.ENDPOINT.CREATE.IN.URI", in.Uri),
	)
}

type EndpointUpdateOut struct {
	*entities.Endpoint
}
