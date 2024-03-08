package usecase

import (
	"context"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"gorm.io/gorm"
)

type Workspace interface {
	Create(ctx context.Context, in *WorkspaceCreateIn) (*WorkspaceCreateOut, error)
	Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error)
	Update(ctx context.Context, in *WorkspaceUpdateIn) (*WorkspaceUpdateOut, error)
	Delete(ctx context.Context, in *WorkspaceDeleteIn) (*WorkspaceDeleteOut, error)
}

type workspace struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	orm    *gorm.DB
}
