package db

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/entities"
)

type Workspace interface {
	Get(ctx context.Context, id string) (*entities.Workspace, error)
}

type WorkspaceCredentials interface {
	Get(ctx context.Context, id string) (*entities.WorkspaceCredentials, error)
}
