package ds

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) Scan(ctx context.Context, epId string, query *datastore.ScanningQuery) chan *datastore.ScanningRecord[[]entities.Request] {
	ch := make(chan *datastore.ScanningRecord[[]entities.Request], 1)
	go sql.scan(ctx, epId, query, ch)
	return ch
}

func (sql *SqlRequest) scan(ctx context.Context, epId string, query *datastore.ScanningQuery, ch chan *datastore.ScanningRecord[[]entities.Request]) {
	defer close(ch)

	// cache the cursor here to use it next round
	cursor := query.Cursor
	idcursor := ""
	for {
		if ctx.Err() != nil {
			return
		}

		tx := sql.client.Model(&entities.Request{}).
			Where("ep_id = ?", epId).
			// the primary key is combined from ep_id, msg_id and id, so to let the database use primary for scanning
			// we need to order by the column ep_id DESC, msg_id DESC first
			// then inside .Sqlx function, the column msg_id DESC will be used to order
			Order("ep_id DESC")

		scanQuery := query.Clone()
		// use previous cursor
		scanQuery.Cursor = cursor

		tx = scanQuery.Sqlx(
			tx,
			&datastore.ScanningCondition{
				PrimaryKeyNs: entities.IdNsEp,
				// we cannot use the column id here
				// because the column msg_id is the next column after the column ep_id in the primary key
				// order by msg_id DESC will be made here
				PrimaryKeyCol: "msg_id",
			},
		)
		// then we need to order by the column id and add the cursor condition
		tx = tx.Order("id DESC")
		if idcursor != "" {
			tx = tx.Where("id < ?", idcursor)
		}

		var data []entities.Request
		if tx := tx.Find(&data); tx.Error != nil {
			ch <- &datastore.ScanningRecord[[]entities.Request]{Error: tx.Error}
			return
		}

		ch <- &datastore.ScanningRecord[[]entities.Request]{Data: data}

		if len(data) < query.Size {
			return
		}

		// refresh cursor for next round
		// by default datastore.ScanningOrderDesc will be used, so the last item must be use as cursor
		cursor = data[len(data)-1].MsgId
		idcursor = data[len(data)-1].Id
	}
}
