package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrRouteGet = errors.New("SDK.ROUTE.GET.ERROR")

func (uc *route) Get(ctx context.Context, in *RouteGetIn) (*RouteGetOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	joinstm := fmt.Sprintf("JOIN %s ON %s.id = %s.ep_id", entities.TableEp, entities.TableEp, entities.TableRt)
	wherestm := fmt.Sprintf("%s.id = ? AND %s.id = ?", entities.TableEp, entities.TableRt)
	selectstm := fmt.Sprintf("%s.*", entities.TableRt)

	doc := &entities.Route{}

	err := uc.orm.
		InnerJoins(joinstm).
		Where(wherestm, in.EpId, in.Id).
		Select(selectstm).
		First(doc).Error
	if err != nil {
		uc.logger.Errorw(ErrRouteGet.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrRouteGet
	}

	out := &RouteGetOut{doc}
	return out, nil
}

type RouteGetIn struct {
	EpId string
	Id   string
}

func (in *RouteGetIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ROUTE.GET.IN.EP_ID", in.EpId, entities.IdNsEp),
		validator.StringStartsWith("SDK.ROUTE.GET.IN.ID", in.Id, entities.IdNsRt),
	)
}

type RouteGetOut struct {
	*entities.Route
}
