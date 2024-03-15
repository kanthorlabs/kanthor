package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrEndpointGetOwn = errors.New("SDK.ENDPOINT.GET.ERROR")

func (uc *endpoint) GetOwn(ctx context.Context, in *EndpointGetOwnIn) (*EndpointGetOwnOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	joinstm := fmt.Sprintf("JOIN %s ON %s.id = %s.app_id", entities.TableApp, entities.TableApp, entities.TableEp)
	wherestm := fmt.Sprintf("%s.ws_id = ? AND %s.id = ?", entities.TableApp, entities.TableEp)
	selectstm := fmt.Sprintf("%s.*", entities.TableEp)

	doc := &entities.Endpoint{}
	err := uc.orm.WithContext(ctx).
		InnerJoins(joinstm).
		Where(wherestm, in.WsId, in.Id).
		Select(selectstm).
		First(doc).Error
	if err != nil {
		uc.logger.Errorw(ErrEndpointGetOwn.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrEndpointGetOwn
	}

	out := &EndpointGetOwnOut{doc}
	return out, nil
}

type EndpointGetOwnIn struct {
	WsId string
	Id   string
}

func (in *EndpointGetOwnIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ENDPOINT.CREATE.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.ENDPOINT.GET.IN.ID", in.Id, entities.IdNsEp),
	)
}

type EndpointGetOwnOut struct {
	*entities.Endpoint
}
