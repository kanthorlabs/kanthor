package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/services/portal/permissions"
)

var ErrCredentialsList = errors.New("PORTAL.CREDENTIALS.LIST.ERROR")

func (uc *credentials) List(ctx context.Context, in *CredentialsListIn) (*CredentialsListOut, error) {
	strategy, err := uc.infra.Passport().Strategy(permissions.Sdk)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsList.Error(), "error", err.Error(), "passport_strategy", permissions.Sdk)
		return nil, ErrCredentialsList
	}

	if err := in.Validate(); err != nil {
		return nil, err
	}

	out := &CredentialsListOut{Data: []*CredentialsAccount{}}

	users, err := uc.infra.Gatekeeper().Users(ctx, in.Tenant)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrCredentialsList
	}

	usernames, maps := cagkmap(users)
	if len(usernames) == 0 {
		return out, nil
	}

	accounts, err := strategy.List(ctx, usernames)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrCredentialsList
	}
	if len(accounts) == 0 {
		return out, nil
	}

	out.Data = append(out.Data, cappmap(maps, accounts)...)
	return out, nil
}

type CredentialsListIn struct {
	Tenant string
}

func (in *CredentialsListIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("PORTAl.CREDENTIALS.LIST.IN.TENANT", in.Tenant),
	)
}

type CredentialsListOut struct {
	Data []*CredentialsAccount
}
