package usecase

import (
	"context"
	"errors"

	gkentities "github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/permissions"
	"gorm.io/gorm"
)

var ErrWorkspaceCreate = errors.New("PORTAL.WORKSPACE.CREATE.ERROR")

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
	doc.SetAuditFacttor(uc.watch.Now())

	err := uc.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(doc).Error; err != nil {
			uc.logger.Errorw(ErrWorkspaceCreate.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrWorkspaceCreate
		}

		evaluation := &gkentities.Evaluation{
			Tenant:   doc.Id,
			Username: doc.OwnerId,
			Role:     permissions.Owner,
		}
		if err := uc.infra.Gatekeeper().Grant(ctx, evaluation); err != nil {
			uc.logger.Errorw(ErrWorkspaceCreate.Error(), "error", err.Error(), "workspace", utils.Stringify(doc), "evaluation", utils.Stringify(evaluation))

			return err
		}

		return nil
	})
	if err != nil {
		// the error is already logged in the transaction
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
