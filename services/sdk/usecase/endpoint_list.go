package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/internal/database/scopes"
)

var ErrEndpointList = errors.New("SDK.ENDPOINT.LIST.ERROR")

func (uc *endpoint) List(ctx context.Context, in *EndpointListIn) (*EndpointListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	model := &entities.Endpoint{}
	base := uc.orm.WithContext(ctx).Model(model).Scopes(scopes.UseEp(in.WsId))
	if in.AppId != "" {
		base = base.Where(fmt.Sprintf("%s.app_id = ?", entities.TableEp), in.AppId)
	}

	var count int64
	var docs []*entities.Endpoint

	if err := in.Query.SqlxCount(base, model.PrimaryProp(), model.SearchProps()).Count(&count).Error; err != nil {
		uc.logger.Errorw(ErrEndpointList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrEndpointList
	}
	if err := in.Query.Sqlx(base, model.PrimaryProp(), model.SearchProps()).Find(&docs).Error; err != nil {
		uc.logger.Errorw(ErrEndpointList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrEndpointList
	}

	return &EndpointListOut{Count: count, Data: docs}, nil
}

type EndpointListIn struct {
	WsId  string
	AppId string
	Query *database.PagingQuery
}

func (in *EndpointListIn) Validate() error {
	err := validator.Validate(
		validator.StringStartsWith("SDK.ENDPOINT.LIST.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWithIfNotEmpty("SDK.ENDPOINT.LIST.IN.APP_ID", in.AppId, entities.IdNsApp),
		validator.PointerNotNil("SDK.ENDPOINT.LIST.IN.QUERY", in.Query),
	)
	if err != nil {
		return err
	}

	if err := in.Query.Validate(); err != nil {
		return err
	}
	return nil
}

type EndpointListOut struct {
	Count int64
	Data  []*entities.Endpoint
}
