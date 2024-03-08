package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
)

var ErrWorkspaceGet = errors.New("PORTAL.WORKSPACE.GET.ERROR")

func (uc *workspace) Get(ctx context.Context, in *WorkspaceGetIn) (*WorkspaceGetOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Workspace{}
	err := uc.orm.Where("id = ?", in.Id).First(doc).Error
	if err != nil {
		uc.logger.Errorw(ErrWorkspaceGet.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrWorkspaceGet
	}

	out := &WorkspaceGetOut{doc}
	return out, nil
}

type WorkspaceGetIn struct {
	Id string
}

func (in *WorkspaceGetIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("PORTAl.WORKSPACE.GET.IN.ID", in.Id, entities.IdNsWs),
	)
}

type WorkspaceGetOut struct {
	*entities.Workspace
}
