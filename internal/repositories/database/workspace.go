package database

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/repositories/database/entities"
	"gorm.io/gorm"
)

type Workspace interface {
	Create(ctx context.Context, doc *entities.Workspace) error
	// Get(ctx context.Context, id string) (*entities.Workspace, error)
	// ListOwned(ctx context.Context, owner string) ([]entities.Workspace, error)
	// ListByIds(ctx context.Context, ids []string) ([]entities.Workspace, error)
	// Update(ctx context.Context, doc *entities.Workspace) (*entities.Workspace, error)
	// Delete(ctx context.Context, id string) (*entities.Workspace, error)
}

type sqlxws struct {
	client *gorm.DB
}

func (db *sqlxws) Create(ctx context.Context, doc *entities.Workspace) error {
	if err := doc.Validate(); err != nil {
		return err
	}

	if err := db.client.Create(doc).Error; err != nil {
		return err
	}
	return nil
}
