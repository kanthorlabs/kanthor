package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrEndpointList = errors.New("SDK.ENDPOINT.LIST.ERROR")

func (uc *endpoint) List(ctx context.Context, in *EndpointListIn) (*EndpointListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	joinstm := fmt.Sprintf("JOIN %s ON %s.id = %s.app_id", entities.TableApp, entities.TableApp, entities.TableEp)
	wherestm := fmt.Sprintf("%s.id = ?", entities.TableApp)
	selectstm := fmt.Sprintf("%s.*", entities.TableEp)

	var count int64
	var docs []*entities.Endpoint

	model := &entities.Endpoint{}
	base := uc.orm.
		InnerJoins(joinstm).
		Where(wherestm, in.AppId).
		Select(selectstm)

	if err := in.Query.SqlxCount(base, model.PrimaryProp(), model.SearchProps()).Count(&count).Error; err != nil {
		return nil, ErrEndpointList
	}
	if err := in.Query.Sqlx(base, model.PrimaryProp(), model.SearchProps()).Find(&docs).Error; err != nil {
		return nil, ErrEndpointList
	}

	return &EndpointListOut{Count: count, Data: docs}, nil
}

type EndpointListIn struct {
	AppId string
	Query *database.PagingQuery
}

func (in *EndpointListIn) Validate() error {
	err := validator.Validate(
		validator.StringStartsWith("SDK.ENDPOINT.LIST.IN.APP_ID", in.AppId, entities.IdNsApp),
	)
	if err != nil {
		return err
	}

	if in.Query != nil {
		return in.Query.Validate()
	}
	// must clone it to get the new one
	in.Query = database.DefaultPagingQuery.Clone()
	return nil
}

type EndpointListOut struct {
	Count int64
	Data  []*entities.Endpoint
}