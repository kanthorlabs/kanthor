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

var ErrWorkspaceDelete = errors.New("PORTAL.WORKSPACE.DELETE.ERROR")

func (uc *workspace) Delete(ctx context.Context, in *WorkspaceDeleteIn) (*WorkspaceDeleteOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	doc := &entities.Workspace{}

	err := uc.orm.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Where("id = ?", in.Id).First(doc).Error
		if err != nil {
			uc.logger.Errorw(ErrWorkspaceDelete.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return ErrWorkspaceDelete
		}

		doc.SetAuditFacttor(in.Modifier, uc.watch.Now())
		if err := tx.Save(doc).Error; err != nil {
			uc.logger.Errorw(ErrWorkspaceDelete.Error(), "error", err.Error(), "in", utils.Stringify(in), "workspace", utils.Stringify(doc))
			return ErrWorkspaceDelete
		}

		if err := tx.Delete(doc).Error; err != nil {
			uc.logger.Errorw(ErrWorkspaceDelete.Error(), "error", err.Error(), "in", utils.Stringify(in), "workspace", utils.Stringify(doc))
			return ErrWorkspaceDelete
		}

		return nil
	})
	if err != nil {
		// the error is already logged in the transaction
		return nil, err
	}

	out := &WorkspaceDeleteOut{doc}
	return out, nil
}

type WorkspaceDeleteIn struct {
	Modifier string
	Id       string
}

func (in *WorkspaceDeleteIn) Validate() error {
	return validator.Validate(
		validator.StringRequired("PORTAl.WORKSPACE.DELETE.IN.MODIFIER", in.Modifier),
		validator.StringStartsWith("PORTAl.WORKSPACE.DELETE.IN.ID", in.Id, entities.IdNsWs),
	)
}

type WorkspaceDeleteOut struct {
	*entities.Workspace
}
