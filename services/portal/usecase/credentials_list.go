package usecase

import (
	"context"
	"errors"
	"slices"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/services/portal/permissions"
)

var ErrCredentialsList = errors.New("PORTAL.CREDENTIALS.LIST.ERROR")

func (uc *credentials) List(ctx context.Context, in *CredentialsListIn) (*CredentialsListOut, error) {
	strategy, err := uc.infra.Passport().Strategy(permissions.Sdk)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsCreate.Error(), "error", err.Error(), "passport_strategy", permissions.Sdk)
		return nil, ErrCredentialsCreate
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

	maps := map[string]*CredentialsAccount{}
	usernames := []string{}
	for i := range users {
		isSdk := slices.Contains(users[i].Roles, permissions.Sdk)
		if isSdk {
			maps[users[i].Username] = &CredentialsAccount{Username: users[i].Username, Roles: users[i].Roles}
			usernames = append(usernames, users[i].Username)
		}
	}

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

	for i := range accounts {
		if account, ok := maps[accounts[i].Username]; ok {
			account.Name = accounts[i].Name
			account.Metadata = accounts[i].Metadata
			account.CreatedAt = accounts[i].CreatedAt
			account.UpdatedAt = accounts[i].UpdatedAt
			out.Data = append(out.Data, account)
		}
	}

	return out, nil
}

type CredentialsListIn struct {
	Tenant string
}

func (in *CredentialsListIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("PORTAl.CREDENTIALS.CREATE.IN.TENANT", in.Tenant),
	)
}

type CredentialsListOut struct {
	Data []*CredentialsAccount
}
