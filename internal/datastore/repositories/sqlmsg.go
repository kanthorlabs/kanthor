package repositories

import (
	"context"
	"fmt"
	"strings"

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

func (repo *sqlmsg) Get(ctx context.Context, pks []entities.MessagePk) ([]*entities.Message, error) {
	var messages []*entities.Message

	if len(pks) == 0 {
		return messages, nil
	}

	wheres := make([]string, len(pks))
	values := map[string]interface{}{}

	for i := range pks {
		wheres[i] = fmt.Sprintf("(app_id = @app_id_%d AND id = @id_%d)", i, i)
		values[fmt.Sprintf("app_id_%d", i)] = pks[i].AppId
		values[fmt.Sprintf("id_%d", i)] = pks[i].Id
	}

	wherestm := strings.Join(wheres, " OR ")
	err := repo.client.Model(&entities.Message{}).Where(wherestm, values).Find(&messages).Error

	return messages, err
}
