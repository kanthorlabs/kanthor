package usecase

import (
	"context"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"gorm.io/gorm"
)

type Credentials interface {
	Create(ctx context.Context, in *CredentialsCreateIn) (*CredentialsCreateOut, error)
	List(ctx context.Context, in *CredentialsListIn) (*CredentialsListOut, error)
	// Expire(ctx context.Context, in *CredentialsExpireIn) (*CredentialsExpireOut, error)
}

type credentials struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	orm    *gorm.DB
}

type CredentialsAccount struct {
	Username  string
	Roles     []string
	Name      string
	Metadata  *safe.Metadata
	CreatedAt int64
	UpdatedAt int64
}
