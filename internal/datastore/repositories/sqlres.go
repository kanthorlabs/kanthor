package repositories

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type sqlres struct {
	client *gorm.DB
}

func (repo *sqlres) Save(ctx context.Context, docs []*entities.Response) error {
	return repo.client.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&docs).Error
}
