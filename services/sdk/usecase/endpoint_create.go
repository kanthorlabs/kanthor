package usecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var SecretLength = 64
var ErrEndpointCreate = errors.New("SDK.ENDPOINT.CREATE.ERROR")

func (uc *endpoint) Create(ctx context.Context, in *EndpointCreateIn) (*EndpointCreateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Endpoint{
		AppId:  in.AppId,
		Name:   in.Name,
		Method: in.Method,
		Uri:    in.Uri,
	}
	doc.SetId()
	doc.SetAuditFacttor(in.Modifier, uc.watch.Now())
	doc.SecretKey = utils.RandomString(SecretLength)

	if err := uc.orm.Create(doc).Error; err != nil {
		uc.logger.Errorw(ErrEndpointCreate.Error(), "error", err.Error(), "in", utils.Stringify(in), "endpoint", utils.Stringify(doc))
		return nil, ErrEndpointCreate
	}

	out := &EndpointCreateOut{doc}
	return out, nil
}

type EndpointCreateIn struct {
	Modifier string
	AppId    string
	Name     string
	Method   string
	Uri      string
}

func (in *EndpointCreateIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("SDK.ENDPOINT.CREATE.IN.MODIFIER", in.Modifier),
		validator.StringStartsWith("SDK.ENDPOINT.CREATE.IN.APP_ID", in.AppId, entities.IdNsApp),
		validator.StringRequired("SDK.ENDPOINT.CREATE.IN.NAME", in.Name),
		validator.StringOneOf("SDK.ENDPOINT.CREATE.IN.METHOD", in.Method, []string{http.MethodPost, http.MethodPut}),
		validator.StringUri("SDK.ENDPOINT.CREATE.IN.URI", in.Uri),
	)
}

type EndpointCreateOut struct {
	*entities.Endpoint
}