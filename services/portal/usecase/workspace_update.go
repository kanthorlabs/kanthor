package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrWorkspaceUpdate = errors.New("PORTAL.WORKSPACE.UPDATE.ERROR")

func (uc *workspace) Update(ctx context.Context, in *WorkspaceUpdateIn) (*WorkspaceUpdateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Workspace{}

	err := uc.orm.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthShare}).Where("id = ?", in.Id).First(doc).Error
		if err != nil {
			uc.logger.Errorw(ErrWorkspaceUpdate.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrWorkspaceUpdate
		}

		doc.Name = in.Name
		doc.SetAuditFacttor(uc.watch.Now(), in.Modifier)
		if err := tx.Save(doc).Error; err != nil {
			uc.logger.Errorw(ErrWorkspaceUpdate.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrWorkspaceUpdate
		}

		return nil
	})
	if err != nil {
		// the error is already logged in the transaction
		return nil, err
	}

	out := &WorkspaceUpdateOut{doc}
	return out, nil
}

type WorkspaceUpdateIn struct {
	Modifier string
	Id       string
	Name     string
}

func (in *WorkspaceUpdateIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("PORTAl.WORKSPACE.UPDATE.IN.MODIFIER", in.Modifier),
		validator.StringStartsWith("PORTAl.WORKSPACE.UPDATE.IN.ID", in.Id, entities.IdNsWs),
		validator.StringRequired("PORTAl.WORKSPACE.UPDATE.IN.NAME", in.Name),
	)
}

type WorkspaceUpdateOut struct {
	*entities.Workspace
}
