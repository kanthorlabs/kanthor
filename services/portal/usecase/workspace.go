package usecase

import (
	"context"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/repositories"
)

type Workspace interface {
	Create(ctx context.Context, in *WorkspaceCreateIn) (*WorkspaceCreateOut, error)
	Update(ctx context.Context, in *WorkspaceUpdateIn) (*WorkspaceUpdateOut, error)
	List(ctx context.Context, in *WorkspaceListIn) (*WorkspaceListOut, error)
	Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error)

	Export(ctx context.Context, in *WorkspaceExportIn) (*WorkspaceExportOut, error)
	Import(ctx context.Context, in *WorkspaceImportIn) (*WorkspaceImportOut, error)
}

type workspace struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
