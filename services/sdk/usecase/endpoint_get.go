package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/internal/database/scopes"
)

var ErrEndpointGet = errors.New("SDK.ENDPOINT.GET.ERROR")

func (uc *endpoint) Get(ctx context.Context, in *EndpointGetIn) (*EndpointGetOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Endpoint{}
	err := uc.orm.WithContext(ctx).
		Scopes(scopes.UseEp(in.WsId)).
		Where(fmt.Sprintf("%s.id = ?", entities.TableEp), in.Id).
		First(doc).Error
	if err != nil {
		uc.logger.Errorw(ErrEndpointGet.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrEndpointGet
	}

	out := &EndpointGetOut{doc}
	return out, nil
}

type EndpointGetIn struct {
	WsId string
	Id   string
}

func (in *EndpointGetIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ENDPOINT.GET.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.ENDPOINT.GET.IN.ID", in.Id, entities.IdNsEp),
	)
}

type EndpointGetOut struct {
	*entities.Endpoint
}
