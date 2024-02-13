package db

import (
	"context"
	"fmt"

	"github.com/kanthorlabs/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlWorkspace struct {
	client *gorm.DB
}

func (sql *SqlWorkspace) Get(ctx context.Context, id string) (*entities.Workspace, error) {
	doc := &entities.Workspace{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
