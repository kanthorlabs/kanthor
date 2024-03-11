package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/services/portal/permissions"
)

var ErrCredentialsExpire = errors.New("PORTAL.CREDENTIALS.EXPIRE.ERROR")

func (uc *credentials) Expire(ctx context.Context, in *CredentialsExpireIn) (*CredentialsExpireOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	accounts, err := uc.get(ctx, in.Tenant, in.Username)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsExpire.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrCredentialsExpire
	}
	if len(accounts) == 0 {
		uc.logger.Errorw(ErrCredentialsExpire.Error(), "error", "no SDK user details was found", "in", utils.Stringify(in), "accounts", utils.Stringify(accounts))
		return nil, ErrCredentialsExpire
	}

	out := &CredentialsExpireOut{accounts[0]}
	out.DeactivatedAt = uc.watch.Now().Add(time.Millisecond * time.Duration(in.ExpiresIn)).UnixMilli()

	// the strategy is already checked in the .get method
	strategy, _ := uc.infra.Passport().Strategy(permissions.Sdk)
	if err := strategy.Deactivate(ctx, in.Username, out.DeactivatedAt); err != nil {
		uc.logger.Errorw(ErrCredentialsExpire.Error(), "error", err.Error(), "in", utils.Stringify(in), "account", utils.Stringify(out))
		return nil, ErrCredentialsExpire
	}

	return out, nil
}

type CredentialsExpireIn struct {
	Tenant    string
	Username  string
	ExpiresIn int64
}

func (in *CredentialsExpireIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("PORTAl.CREDENTIALS.EXPIRE.IN.TENANT", in.Tenant),
		validator.StringRequired("PORTAl.CREDENTIALS.EXPIRE.IN.USERNAME", in.Username),
		// 30 minutes to 7 days
		validator.NumberInRange("PORTAl.CREDENTIALS.EXPIRE.IN.EXPIRES_IN", in.ExpiresIn, 1800000, 604800000),
	)
}

type CredentialsExpireOut struct {
	*CredentialsAccount
}
