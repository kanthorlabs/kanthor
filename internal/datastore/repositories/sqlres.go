package repositories

import (
	"context"
	"fmt"
	"strings"

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

func (repo *sqlres) Get(ctx context.Context, pks []entities.ResponsePk) ([]*entities.Response, error) {
	var responses []*entities.Response

	if len(pks) == 0 {
		return responses, nil
	}

	wheres := make([]string, len(pks))
	values := map[string]interface{}{}

	for i := range pks {
		wheres[i] = fmt.Sprintf("(ep_id = @ep_id_%d AND msg_id = @msg_id_%d)", i, i)
		values[fmt.Sprintf("ep_id_%d", i)] = pks[i].EpId
		values[fmt.Sprintf("msg_id_%d", i)] = pks[i].MsgId
	}

	wherestm := strings.Join(wheres, " OR ")
	err := repo.client.Model(&entities.Response{}).Where(wherestm, values).Find(&responses).Error

	return responses, err
}
