package usecase

import (
	"context"

	gkentities "github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/repositories/database/entities"
	"github.com/kanthorlabs/kanthor/services/portal/permissions"
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

	// @TODO: add transaction
	if err := uc.repos.Workspace().Create(ctx, doc); err != nil {
		return nil, err
	}

	evaluation := &gkentities.Evaluation{
		Tenant:   doc.Id,
		Username: doc.OwnerId,
		Role:     permissions.RoleOwner,
	}
	if err := uc.infra.Gatekeeper().Grant(ctx, evaluation); err != nil {
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
