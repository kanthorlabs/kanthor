package ds

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlMessage struct {
	client *gorm.DB
}

func (sql *SqlMessage) Scan(ctx context.Context, appId string, query *datastore.ScanningQuery) chan *datastore.ScanningRecord[[]entities.Message] {
	ch := make(chan *datastore.ScanningRecord[[]entities.Message], 1)
	go sql.scan(ctx, appId, query, ch)
	return ch
}

func (sql *SqlMessage) scan(ctx context.Context, appId string, query *datastore.ScanningQuery, ch chan *datastore.ScanningRecord[[]entities.Message]) {
	defer close(ch)

	// cache the cursor here to use it next round
	cursor := query.Cursor
	for {
		if ctx.Err() != nil {
			return
		}

		tx := sql.client.Model(&entities.Message{}).
			Where("app_id = ?", appId).
			// the primary key is combined from app_id and id, so to let the database use primary for scanning
			// we need to order by the column app_id DESC first
			// then inside .Sqlx function, the column id DESC will be used to order
			Order("app_id DESC")

		scanQuery := query.Clone()
		// use previous cursor
		scanQuery.Cursor = cursor

		tx = scanQuery.Sqlx(
			tx,
			&datastore.ScanningCondition{
				PrimaryKeyNs:  entities.IdNsEp,
				PrimaryKeyCol: "id",
			},
		)
		var data []entities.Message
		if tx := tx.Find(&data); tx.Error != nil {
			ch <- &datastore.ScanningRecord[[]entities.Message]{Error: tx.Error}
			return
		}

		ch <- &datastore.ScanningRecord[[]entities.Message]{Data: data}

		if len(data) < query.Size {
			return
		}

		// refresh cursor for next round
		// by default datastore.ScanningOrderDesc will be used, so the last item must be use as cursor
		cursor = data[len(data)-1].Id
	}
}
