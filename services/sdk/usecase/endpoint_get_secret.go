package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/cipher/encryption"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/internal/database/scopes"
	"gorm.io/gorm/clause"
)

var ErrEndpointGetSecret = errors.New("SDK.ENDPOINT.GET_SECRET.ERROR")

func (uc *endpoint) GetSecret(ctx context.Context, in *EndpointGetSecretIn) (*EndpointGetSecretOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Endpoint{}
	err := uc.orm.WithContext(ctx).
		Clauses(clause.Locking{Strength: clause.LockingStrengthShare}).
		Scopes(scopes.UseEp(in.WsId)).
		Where(fmt.Sprintf("%s.id = ?", entities.TableEp), in.Id).
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
	WsId string
	Id   string
}

func (in *EndpointGetSecretIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.ENDPOINT.GET.IN.WS_ID", in.WsId, entities.IdNsWs),
		validator.StringStartsWith("SDK.ENDPOINT.GET_SECRET.IN.ID", in.Id, entities.IdNsEp),
	)
}

type EndpointGetSecretOut struct {
	*entities.Endpoint
	DescryptedSecretKey string
}
