package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type WorkspaceCredentialsUpdateIn struct {
	WsId      string
	Id        string
	Name      string
	ExpiredAt int64
}

func (req *WorkspaceCredentialsUpdateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("ws_id", req.WsId, entities.IdNsWs),
		validator.StringStartsWith("id", req.Id, entities.IdNsWsc),
		validator.StringRequired("name", req.Name),
	)
}

type WorkspaceCredentialsUpdateOut struct {
	Doc *entities.WorkspaceCredentials
}

func (uc *workspaceCredentials) Update(ctx context.Context, in *WorkspaceCredentialsUpdateIn) (*WorkspaceCredentialsUpdateOut, error) {
	wsc, err := uc.repositories.Database().WorkspaceCredentials().Get(ctx, in.WsId, in.Id)
	if err != nil {
		return nil, err
	}

	wsc.Name = in.Name
	wsc.ExpiredAt = in.ExpiredAt
	wsc.SetAT(uc.infra.Timer.Now())

	doc, err := uc.repositories.Database().WorkspaceCredentials().Update(ctx, wsc)
	if err != nil {
		return nil, err
	}

	out := &WorkspaceCredentialsUpdateOut{Doc: doc}
	return out, nil
}
