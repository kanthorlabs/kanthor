package db

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlEndpoint struct {
	client *gorm.DB
}

func (sql *SqlEndpoint) Scan(ctx context.Context, query *datastore.ScanningQuery) chan *datastore.ScanningRecord[[]entities.Endpoint] {
	ch := make(chan *datastore.ScanningRecord[[]entities.Endpoint], 1)
	go sql.scan(ctx, query, ch)
	return ch
}

func (sql *SqlEndpoint) scan(ctx context.Context, query *datastore.ScanningQuery, ch chan *datastore.ScanningRecord[[]entities.Endpoint]) {
	defer close(ch)

	// cache the cursor here to use it next round
	cursor := query.Cursor
	for {
		if ctx.Err() != nil {
			return
		}

		scanQuery := query.Clone()
		// use previous cursor
		scanQuery.Cursor = cursor

		tx := scanQuery.Sqlx(
			sql.client.Model(&entities.Endpoint{}),
			&datastore.ScanningCondition{
				PrimaryKeyNs:  entities.IdNsEp,
				PrimaryKeyCol: "id",
			},
		)

		var data []entities.Endpoint
		if tx := tx.Find(&data); tx.Error != nil {
			ch <- &datastore.ScanningRecord[[]entities.Endpoint]{Error: tx.Error}
			return
		}

		ch <- &datastore.ScanningRecord[[]entities.Endpoint]{Data: data}

		if len(data) < query.Size {
			return
		}

		// refresh cursor for next round
		// by default datastore.ScanningOrderDesc will be used, so the last item must be use as cursor
		cursor = data[len(data)-1].Id
	}
}
