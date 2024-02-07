package usecase

import (
	"context"

	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/repositories"
)

type WorkspaceCredentials interface {
	Generate(ctx context.Context, in *WorkspaceCredentialsGenerateIn) (*WorkspaceCredentialsGenerateOut, error)
	Update(ctx context.Context, in *WorkspaceCredentialsUpdateIn) (*WorkspaceCredentialsUpdateOut, error)
	Expire(ctx context.Context, in *WorkspaceCredentialsExpireIn) (*WorkspaceCredentialsExpireOut, error)
	List(ctx context.Context, in *WorkspaceCredentialsListIn) (*WorkspaceCredentialsListOut, error)
	Get(ctx context.Context, in *WorkspaceCredentialsGetIn) (*WorkspaceCredentialsGetOut, error)
}

type workspaceCredentials struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
