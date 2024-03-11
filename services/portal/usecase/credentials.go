package usecase

import (
	"context"
	"slices"

	"github.com/kanthorlabs/common/clock"
	gkentities "github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/logging"
	ppentities "github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/permissions"
	"gorm.io/gorm"
)

type Credentials interface {
	Create(ctx context.Context, in *CredentialsCreateIn) (*CredentialsCreateOut, error)
	List(ctx context.Context, in *CredentialsListIn) (*CredentialsListOut, error)
	Get(ctx context.Context, in *CredentialsGetIn) (*CredentialsGetOut, error)
	Expire(ctx context.Context, in *CredentialsExpireIn) (*CredentialsExpireOut, error)
}

type credentials struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	orm    *gorm.DB
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

func cagkmap(users []gkentities.User) ([]string, map[string]*CredentialsAccount) {
	usernames := []string{}
	maps := map[string]*CredentialsAccount{}
	for i := range users {
		// the credentials are only for the SDK
		isSdk := slices.Contains(users[i].Roles, permissions.Sdk)
		if isSdk {
			usernames = append(usernames, users[i].Username)
			maps[users[i].Username] = &CredentialsAccount{
				Username: users[i].Username,
				Roles:    users[i].Roles,
			}
		}
	}
	return usernames, maps
}

func cappmap(maps map[string]*CredentialsAccount, accounts []*ppentities.Account) []*CredentialsAccount {
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
