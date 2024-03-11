package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
)

var ErrCredentialsGet = errors.New("PORTAL.CREDENTIALS.GET.ERROR")

func (uc *credentials) Get(ctx context.Context, in *CredentialsGetIn) (*CredentialsGetOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	accounts, err := uc.get(ctx, in.Tenant, in.Username)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsGet.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrCredentialsGet
	}
	if len(accounts) == 0 {
		uc.logger.Errorw(ErrCredentialsGet.Error(), "error", "no SDK user details was found", "in", utils.Stringify(in), "accounts", utils.Stringify(accounts))
		return nil, ErrCredentialsGet
	}

	out := &CredentialsGetOut{accounts[0]}
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
