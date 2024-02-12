package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

type WorkspaceGetIn struct {
	AccId string
	Id    string
}

func (in *WorkspaceGetIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("acc_id", in.AccId),
		validator.StringStartsWith("id", in.Id, entities.IdNsWs),
	)
}

type WorkspaceGetOut struct {
	Doc *entities.Workspace
}

func (uc *workspace) Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error) {
	ws, err := uc.repositories.Database().Workspace().Get(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	isOwner := ws.OwnerId == in.AccId
	if !isOwner {
		return nil, errors.New("SDK.USECASE.WORKSPACE.GET.NOT_OWN.ERROR")
	}

	return &WorkspaceGetOut{ws}, nil
}
