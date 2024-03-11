package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/services/portal/permissions"
)

var ErrCredentialsGet = errors.New("PORTAL.CREDENTIALS.GET.ERROR")

func (uc *credentials) Get(ctx context.Context, in *CredentialsGetIn) (*CredentialsGetOut, error) {
	strategy, err := uc.infra.Passport().Strategy(permissions.Sdk)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsGet.Error(), "error", err.Error(), "passport_strategy", permissions.Sdk)
		return nil, ErrCredentialsGet
	}

	if err := in.Validate(); err != nil {
		return nil, err
	}

	users, err := uc.infra.Gatekeeper().Users(ctx, in.Tenant)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsGet.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrCredentialsGet
	}

	usernames, maps := cagkmap(users)
	if len(usernames) == 0 {
		uc.logger.Errorw(ErrCredentialsGet.Error(), "error", "no SDK user was found", "in", utils.Stringify(in), "users", utils.Stringify(users))
		return nil, ErrCredentialsGet
	}

	accounts, err := strategy.List(ctx, usernames)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsGet.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrCredentialsGet
	}
	if len(accounts) == 0 {
		uc.logger.Errorw(ErrCredentialsGet.Error(), "error", "no SDK user details was found", "in", utils.Stringify(in), "users", utils.Stringify(users))
		return nil, ErrCredentialsGet
	}

	ca := cappmap(maps, accounts)
	out := &CredentialsGetOut{ca[0]}
	return out, nil
}

type CredentialsGetIn struct {
	Tenant   string
	Username string
}

func (in *CredentialsGetIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("PORTAl.CREDENTIALS.GET.IN.TENANT", in.Tenant),
		validator.StringRequired("PORTAl.CREDENTIALS.GET.IN.USERNAME", in.Username),
	)
}

type CredentialsGetOut struct {
	*CredentialsAccount
}
