package usecase

import (
	"context"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type WorkspaceCredentialsListIn struct {
	*database.PagingQuery
	WsId string
}

func (in *WorkspaceCredentialsListIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
	)
}

type WorkspaceCredentialsListOut struct {
	Data  []entities.WorkspaceCredentials
	Count int64
}

func (uc *workspaceCredentials) List(ctx context.Context, in *WorkspaceCredentialsListIn) (*WorkspaceCredentialsListOut, error) {
	data, err := uc.repositories.Database().WorkspaceCredentials().List(ctx, in.WsId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	count, err := uc.repositories.Database().WorkspaceCredentials().Count(ctx, in.WsId, in.PagingQuery)
	if err != nil {
		return nil, err
	}

	out := &WorkspaceCredentialsListOut{Data: data, Count: count}
	return out, nil
}
