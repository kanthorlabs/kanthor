package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type WorkspaceUpdateIn struct {
	AccId string
	Id    string
	Name  string
}

func (in *WorkspaceUpdateIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("acc_id", in.AccId),
		validator.StringStartsWith("id", in.Id, entities.IdNsWs),
		validator.StringRequired("name", in.Name),
	)
}

type WorkspaceUpdateOut struct {
	Doc *entities.Workspace
}

func (uc *workspace) Update(ctx context.Context, in *WorkspaceUpdateIn) (*WorkspaceUpdateOut, error) {
	// the .Get usecase is implemented the logic of authn & authz
	getout, err := uc.Get(ctx, &WorkspaceGetIn{AccId: in.AccId, Id: in.Id})
	if err != nil {
		return nil, err
	}

	ws, err := uc.repositories.Database().Workspace().Get(ctx, getout.Doc.Id)
	if err != nil {
		return nil, err
	}

	ws.Name = in.Name
	ws.SetAT(uc.infra.Timer.Now())
	doc, err := uc.repositories.Database().Workspace().Update(ctx, ws)
	if err != nil {
		return nil, err
	}

	return &WorkspaceUpdateOut{Doc: doc}, nil
}
