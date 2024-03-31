package repositories

import (
	"context"
	"fmt"
	"strings"

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

func (repo *sqlreq) Get(ctx context.Context, pks []entities.RequestPk) ([]*entities.Request, error) {
	var requests []*entities.Request

	if len(pks) == 0 {
		return requests, nil
	}

	wheres := make([]string, len(pks))
	values := map[string]interface{}{}

	for i := range pks {
		wheres[i] = fmt.Sprintf("(ep_id = @ep_id_%d AND msg_id = @msg_id_%d)", i, i)
		values[fmt.Sprintf("ep_id_%d", i)] = pks[i].EpId
		values[fmt.Sprintf("msg_id_%d", i)] = pks[i].MsgId
	}

	wherestm := strings.Join(wheres, " OR ")
	err := repo.client.Model(&entities.Request{}).Where(wherestm, values).Find(&requests).Error

	return requests, err
}
