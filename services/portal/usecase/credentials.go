package usecase

import (
	"context"
	"errors"
	"slices"

	"github.com/kanthorlabs/common/clock"
	gkentities "github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/logging"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/permissions"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"gorm.io/gorm"
)

type Credentials interface {
	Create(ctx context.Context, in *CredentialsCreateIn) (*CredentialsCreateOut, error)
	List(ctx context.Context, in *CredentialsListIn) (*CredentialsListOut, error)
	Get(ctx context.Context, in *CredentialsGetIn) (*CredentialsGetOut, error)
	Update(ctx context.Context, in *CredentialsUpdateIn) (*CredentialsUpdateOut, error)
	Expire(ctx context.Context, in *CredentialsExpireIn) (*CredentialsExpireOut, error)
}

type credentials struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	orm    *gorm.DB
}

func (uc *credentials) get(ctx context.Context, tenant, username string) ([]*CredentialsAccount, error) {
	strategy, err := uc.infra.Passport().Strategy(permissions.Sdk)
	if err != nil {
		return nil, errors.New("PORTAl.CREDENTIALS.PASSPORT.NO_STRATEGY.ERROR")
	}

	users, err := uc.infra.Gatekeeper().Users(ctx, tenant)
	if err != nil {
		return nil, err
	}

	usernames, maps := cagkmap(users, username)
	if len(usernames) == 0 {
		return []*CredentialsAccount{}, nil
	}

	accounts, err := strategy.Management().List(ctx, usernames)
	if err != nil {
		return nil, errors.New("PORTAl.CREDENTIALS.PASSPORT.ERROR")
	}

	return cappmap(accounts, maps), nil
}

type CredentialsAccount struct {
	Username      string
	Roles         []string
	Name          string
	Metadata      *safe.Metadata
	CreatedAt     int64
	UpdatedAt     int64
	DeactivatedAt int64
}

func cagkmap(users []gkentities.User, find string) ([]string, map[string]*CredentialsAccount) {
	usernames := []string{}
	maps := map[string]*CredentialsAccount{}
	for i := range users {
		// the credentials are only for the SDK, ignore the rest type
		//lint:ignore S1002 this is a valid comparison
		othertype := slices.Contains(users[i].Roles, permissions.Sdk) != true
		if othertype {
			continue
		}

		if find != "" && users[i].Username != find {
			continue
		}

		usernames = append(usernames, users[i].Username)
		maps[users[i].Username] = &CredentialsAccount{
			Username: users[i].Username,
			Roles:    users[i].Roles,
		}
	}
	return usernames, maps
}

func cappmap(accounts []*ppentities.Account, maps map[string]*CredentialsAccount) []*CredentialsAccount {
	ca := []*CredentialsAccount{}
	for i := range accounts {
		if account, ok := maps[accounts[i].Username]; ok {
			account.Name = accounts[i].Name
			account.Metadata = accounts[i].Metadata
			account.CreatedAt = accounts[i].CreatedAt
			account.UpdatedAt = accounts[i].UpdatedAt
			account.DeactivatedAt = accounts[i].DeactivatedAt
			ca = append(ca, account)
		}
	}
	return ca
}
