package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
)

var ErrCredentialsList = errors.New("PORTAL.CREDENTIALS.LIST.ERROR")

func (uc *credentials) List(ctx context.Context, in *CredentialsListIn) (*CredentialsListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	accounts, err := uc.get(ctx, in.Tenant, "")
	if err != nil {
		uc.logger.Errorw(ErrCredentialsList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrCredentialsList
	}

	out := &CredentialsListOut{Data: accounts}
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
