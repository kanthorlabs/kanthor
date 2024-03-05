package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/repositories/database/entities"
)

func (uc *workspace) Create(ctx context.Context, in *WorkspaceCreateIn) (*WorkspaceCreateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Workspace{
		OwnerId: in.OwnerId,
		Name:    in.Name,
		Tier:    in.Tier,
	}
	doc.SetId()
	doc.SetAuditTime(uc.watch.Now())

	if err := uc.repos.Workspace().Create(ctx, doc); err != nil {
		return nil, err
	}

	out := &WorkspaceCreateOut{doc}
	return out, nil
}

type WorkspaceCreateIn struct {
	OwnerId string
	Name    string
	Tier    string
}

func (in *WorkspaceCreateIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("PORTAl.WORKSPACE.CREATE.IN.OWNER_ID", in.OwnerId),
		validator.StringRequired("PORTAl.WORKSPACE.CREATE.IN.NAME", in.Name),
		validator.StringRequired("PORTAl.WORKSPACE.CREATE.IN.TIER", in.Tier),
	)
}

type WorkspaceCreateOut struct {
	*entities.Workspace
}
