package repositories

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type sqlreq struct {
	client *gorm.DB
}

func (repo *sqlreq) Save(ctx context.Context, docs []*entities.Request) error {
	return repo.client.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&docs).Error
}
