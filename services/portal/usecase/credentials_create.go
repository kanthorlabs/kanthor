package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/gatekeeper"
	gkentities "github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/idx"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/permissions"
	"go.opentelemetry.io/otel"
)

// bcrypt: password length exceeds 72 bytes
var PasswordLength = 64
var ErrCredentialsCreate = errors.New("PORTAL.CREDENTIALS.CREATE.ERROR")

func (uc *credentials) Create(ctx context.Context, in *CredentialsCreateIn) (*CredentialsCreateOut, error) {
	ctx, span := otel.Tracer("SERVICE.CREDENTIALS.CREATE").Start(ctx, "USECASE")
	defer span.End()

	strategy, err := uc.strategy()
	if err != nil {
		uc.logger.Errorw(ErrCredentialsCreate.Error(), "error", err.Error(), "passport_strategy", permissions.Sdk)
		return nil, ErrCredentialsCreate
	}

	if err := in.Validate(); err != nil {
		return nil, err
	}

	out := &CredentialsCreateOut{
		Tenant:   in.Tenant,
		Username: idx.New(entities.IdNsWsc),
		Password: utils.RandomString(PasswordLength),
	}

	acc := ppentities.Account{
		Username:  out.Username,
		Password:  out.Password,
		Name:      in.Name,
		Metadata:  &safe.Metadata{},
		CreatedAt: uc.watch.Now().UnixMilli(),
		UpdatedAt: uc.watch.Now().UnixMilli(),
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
