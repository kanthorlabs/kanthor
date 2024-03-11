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
	strategy, err := uc.infra.Passport().Strategy(permissions.Sdk)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsExpire.Error(), "error", err.Error(), "passport_strategy", permissions.Sdk)
		return nil, ErrCredentialsExpire
	}

	if err := in.Validate(); err != nil {
		return nil, err
	}

	users, err := uc.infra.Gatekeeper().Users(ctx, in.Tenant)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsExpire.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrCredentialsExpire
	}

	usernames, maps := cagkmap(users)
	if len(usernames) == 0 {
		uc.logger.Errorw(ErrCredentialsExpire.Error(), "error", "no SDK user was found", "in", utils.Stringify(in), "users", utils.Stringify(users))
		return nil, ErrCredentialsExpire
	}

	accounts, err := strategy.List(ctx, usernames)
	if err != nil {
		uc.logger.Errorw(ErrCredentialsExpire.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrCredentialsExpire
	}
	if len(accounts) == 0 {
		uc.logger.Errorw(ErrCredentialsExpire.Error(), "error", "no SDK user details was found", "in", utils.Stringify(in), "users", utils.Stringify(users))
		return nil, ErrCredentialsExpire
	}

	ca := cappmap(maps, accounts)
	out := &CredentialsExpireOut{ca[0]}

	out.DeactivatedAt = uc.watch.Now().Add(time.Millisecond * time.Duration(in.ExpiresIn)).UnixMilli()
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
