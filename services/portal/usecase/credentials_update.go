package usecase

import (
	"context"
	"errors"

	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/services/permissions"
)

var ErrCredentialsUpdate = errors.New("PORTAL.CREDENTIALS.UPDATE.ERROR")

func (uc *credentials) Update(ctx context.Context, in *CredentialsUpdateIn) (*CredentialsUpdateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	accounts, err := uc.get(ctx, in.Tenant, in.Username)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsUpdate.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrCredentialsUpdate
	}
	if len(accounts) == 0 {
		uc.logger.Errorw(ErrCredentialsUpdate.Error(), "error", "no SDK user details was found", "in", utils.Stringify(in), "accounts", utils.Stringify(accounts))
		return nil, ErrCredentialsUpdate
	}

	accounts[0].Name = in.Name
	account := ppentities.Account{
		Username: accounts[0].Username,
		Name:     accounts[0].Name,
	}

	// the strategy is already checked in the .get method
	strategy, _ := uc.infra.Passport().Strategy(permissions.Sdk)
	if err := strategy.Update(ctx, account); err != nil {
		uc.logger.Errorw(ErrCredentialsUpdate.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrCredentialsUpdate
	}

	out := &CredentialsUpdateOut{accounts[0]}
	return out, nil
}

type CredentialsUpdateIn struct {
	Tenant   string
	Username string
	Name     string
}

func (in *CredentialsUpdateIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("PORTAl.CREDENTIALS.UPDATE.IN.TENANT", in.Tenant),
		validator.StringRequired("PORTAl.CREDENTIALS.UPDATE.IN.USERNAME", in.Username),
		validator.StringRequired("PORTAl.CREDENTIALS.UPDATE.IN.NAME", in.Name),
		// 30 minutes to 7 days
	)
}

type CredentialsUpdateOut struct {
	*CredentialsAccount
}
