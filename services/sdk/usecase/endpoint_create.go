package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/kanthorlabs/common/cipher/encryption"
	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/internal/database/scopes"
)

var SecretLength = 128
var ErrEndpointCreate = errors.New("SDK.ENDPOINT.CREATE.ERROR")

func (uc *endpoint) Create(ctx context.Context, in *EndpointCreateIn) (*EndpointCreateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// ensure that the application exists in the requesting workspace before proceeding
	var app entities.Application
	err := uc.orm.WithContext(ctx).Scopes(scopes.UseApp(in.WsId)).
		Where(fmt.Sprintf("%s.id = ?", entities.TableApp), in.AppId).
		First(&app).Error
	if err != nil {
		uc.logger.Errorw(ErrEndpointCreate.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrEndpointCreate
	}

	doc := &entities.Endpoint{
		AppId:  app.Id,
		Name:   in.Name,
		Method: in.Method,
		Uri:    in.Uri,
	}
	doc.SetId()
	doc.SetAuditFacttor(uc.watch.Now())

	secretKey := idx.Build(entities.IdNsEpSec, utils.RandomString(SecretLength))
	secret, err := encryption.Encrypt(uc.conf.Infrastructure.Secrets.Cipher[0], secretKey)
	if err != nil {
		uc.logger.Errorw(ErrEndpointCreate.Error(), "error", err.Error(), "in", utils.Stringify(in), "endpoint", utils.Stringify(doc))
		return nil, ErrEndpointCreate
	}
	doc.SecretKey = secret

	if err := uc.orm.WithContext(ctx).Create(doc).Error; err != nil {
		uc.logger.Errorw(ErrEndpointCreate.Error(), "error", err.Error(), "in", utils.Stringify(in), "endpoint", utils.Stringify(doc))
		return nil, ErrEndpointCreate
	}

	out := &EndpointCreateOut{doc}
	return out, nil
}

type EndpointCreateIn struct {
	WsId   string
	AppId  string
	Name   string
	Method string
	Uri    string
}

func (in *EndpointCreateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ENDPOINT.CREATE.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.ENDPOINT.CREATE.IN.APP_ID", in.AppId, entities.IdNsApp),
		validator.StringRequired("SDK.ENDPOINT.CREATE.IN.NAME", in.Name),
		validator.StringOneOf("SDK.ENDPOINT.CREATE.IN.METHOD", in.Method, []string{http.MethodPost, http.MethodPut}),
		validator.StringUri("SDK.ENDPOINT.CREATE.IN.URI", in.Uri),
	)
}

type EndpointCreateOut struct {
	*entities.Endpoint
}
