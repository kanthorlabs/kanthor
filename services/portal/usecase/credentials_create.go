package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/gatekeeper"
	gkentities "github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/idx"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/services/portal/permissions"
)

var ErrCredentialsCreate = errors.New("PORTAL.CREDENTIALS.CREATE.ERROR")

func (uc *credentials) Create(ctx context.Context, in *CredentialsCreateIn) (*CredentialsCreateOut, error) {
	strategy, err := uc.infra.Passport().Strategy(permissions.Sdk)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsCreate.Error(), "error", err.Error(), "passport_strategy", permissions.Sdk)
		return nil, ErrCredentialsCreate
	}

	if err := in.Validate(); err != nil {
		return nil, err
	}

	out := &CredentialsCreateOut{
		Tenant:   in.Tenant,
		Username: idx.New("sdkacc"),
		Password: utils.RandomString(64),
	}
	hash, err := password.Hash(out.Password)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsCreate.Error(), "error", err.Error(), "in", utils.Stringify(in), "out", utils.Stringify(out))
		return nil, ErrCredentialsCreate
	}

	acc := &ppentities.Account{
		Username:     out.Username,
		PasswordHash: hash,
		Name:         in.Name,
		Metadata:     &safe.Metadata{},
		CreatedAt:    uc.watch.Now().UnixMilli(),
		UpdatedAt:    uc.watch.Now().UnixMilli(),
	}
	// the account must be bound to the tenant
	acc.Metadata.Set(string(gatekeeper.CtxTenantId), out.Tenant)

	if err := strategy.Register(ctx, acc); err != nil {
		uc.logger.Errorw(ErrCredentialsCreate.Error(), "error", err.Error(), "account", utils.Stringify(acc))
		return nil, ErrCredentialsCreate
	}

	evaluation := &gkentities.Evaluation{
		Tenant:   out.Tenant,
		Username: acc.Username,
		Role:     permissions.Sdk,
	}
	if err := uc.infra.Gatekeeper().Grant(ctx, evaluation); err != nil {
		uc.logger.Errorw(ErrCredentialsCreate.Error(), "error", err.Error(), "account", utils.Stringify(acc), "evaluation", utils.Stringify(evaluation))
		return nil, ErrCredentialsCreate
	}

	return out, nil
}

type CredentialsCreateIn struct {
	Tenant string
	Name   string
}

func (in *CredentialsCreateIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("PORTAl.CREDENTIALS.CREATE.IN.TENANT", in.Tenant),
		validator.StringRequired("PORTAl.CREDENTIALS.CREATE.IN.NAME", in.Name),
	)
}

type CredentialsCreateOut struct {
	Tenant   string
	Username string
	Password string
}