package repositories

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type sqlmsg struct {
	client *gorm.DB
}

func (repo *sqlmsg) Save(ctx context.Context, docs []*entities.Message) error {
	return repo.client.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&docs).Error
}
