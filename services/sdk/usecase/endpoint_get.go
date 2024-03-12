package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrEndpointGet = errors.New("SDK.ENDPOINT.GET.ERROR")

func (uc *endpoint) Get(ctx context.Context, in *EndpointGetIn) (*EndpointGetOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	joinstm := fmt.Sprintf("JOIN %s ON %s.id = %s.app_id", entities.TableApp, entities.TableApp, entities.TableEp)
	wherestm := fmt.Sprintf("%s.id = ? AND %s.id = ?", entities.TableApp, entities.TableEp)
	selectstm := fmt.Sprintf("%s.*", entities.TableEp)

	doc := &entities.Endpoint{}

	err := uc.orm.
		InnerJoins(joinstm).
		Where(wherestm, in.AppId, in.Id).
		Select(selectstm).
		First(doc).Error
	if err != nil {
		uc.logger.Errorw(ErrEndpointGet.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrEndpointGet
	}

	out := &EndpointGetOut{doc}
	return out, nil
}

type EndpointGetIn struct {
	AppId string
	Id    string
}

func (in *EndpointGetIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ENDPOINT.CREATE.IN.APP_ID", in.AppId, entities.IdNsApp),
		validator.StringStartsWith("SDK.ENDPOINT.GET.IN.ID", in.Id, entities.IdNsEp),
	)
}

type EndpointGetOut struct {
	*entities.Endpoint
}
