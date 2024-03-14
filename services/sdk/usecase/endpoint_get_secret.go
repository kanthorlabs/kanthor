package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/cipher/encryption"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrEndpointGetSecret = errors.New("SDK.ENDPOINT.GET_SECRET.ERROR")

func (uc *endpoint) GetSecret(ctx context.Context, in *EndpointGetSecretIn) (*EndpointGetSecretOut, error) {
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
		uc.logger.Errorw(ErrEndpointGetSecret.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrEndpointGetSecret
	}

	secretkey, err := encryption.DecryptAny(uc.conf.Infrastructure.Secrets.Cipher, doc.SecretKey)
	if err != nil {
		uc.logger.Errorw(ErrEndpointGetSecret.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrEndpointGetSecret
	}

	out := &EndpointGetSecretOut{doc, secretkey}
	return out, nil
}

type EndpointGetSecretIn struct {
	AppId string
	Id    string
}

func (in *EndpointGetSecretIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ENDPOINT.CREATE.IN.APP_ID", in.AppId, entities.IdNsApp),
		validator.StringStartsWith("SDK.ENDPOINT.GET_SECRET.IN.ID", in.Id, entities.IdNsEp),
	)
}

type EndpointGetSecretOut struct {
	*entities.Endpoint
	DescryptedSecretKey string
}
