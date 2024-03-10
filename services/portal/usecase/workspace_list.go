package usecase

import (
	"context"
	"errors"

	gkentities "github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/samber/lo"
)

var ErrWorkspaceList = errors.New("PORTAL.WORKSPACE.LIST.ERROR")

func (uc *workspace) List(ctx context.Context, in *WorkspaceListIn) (*WorkspaceListOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	tenants, err := uc.infra.Gatekeeper().Tenants(ctx, in.OwnerId)
	if err != nil {
		uc.logger.Errorw(ErrWorkspaceList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrWorkspaceList
	}

	var docs []*entities.Workspace
	// no tenants, return empty list
	if len(tenants) == 0 {
		return &WorkspaceListOut{Data: docs}, nil
	}
	ids := lo.Map(tenants, func(tenant gkentities.Tenant, _ int) string {
		return tenant.Tenant
	})

	err = uc.orm.Where("id IN ?", ids).Find(&docs).Error
	if err != nil {
		uc.logger.Errorw(ErrWorkspaceList.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrWorkspaceList
	}

	return &WorkspaceListOut{Data: docs}, nil
}

type WorkspaceListIn struct {
	OwnerId string
}

func (in *WorkspaceListIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("PORTAl.WORKSPACE.LIST.IN.OWNER_ID", in.OwnerId),
	)
}

type WorkspaceListOut struct {
	Data []*entities.Workspace
}
