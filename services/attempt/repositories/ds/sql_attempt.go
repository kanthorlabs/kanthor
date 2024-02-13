package ds

import (
	"context"
	"fmt"
	"strings"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"github.com/sourcegraph/conc"
	"gorm.io/gorm"
)

type SqlAttempt struct {
	client *gorm.DB
}

func (sql *SqlAttempt) Scan(ctx context.Context, query *datastore.ScanningQuery, next int64, count int) chan *datastore.ScanningRecord[[]entities.Attempt] {
	ch := make(chan *datastore.ScanningRecord[[]entities.Attempt], 1)
	go sql.scan(ctx, query, next, count, ch)
	return ch
}

func (sql *SqlAttempt) scan(ctx context.Context, query *datastore.ScanningQuery, next int64, count int, ch chan *datastore.ScanningRecord[[]entities.Attempt]) {
	defer close(ch)

	// cache the cursor here to use it next round
	cursor := query.Cursor
	for {
		if ctx.Err() != nil {
			return
		}
		tx := sql.client.Model(&entities.Attempt{}).
			Where("completed_at = 0 AND schedule_next <= ? AND schedule_counter < ?", next, count)

		scanQuery := query.Clone()
		// use previous cursor
		scanQuery.Cursor = cursor
		tx = scanQuery.Sqlx(
			tx,
			&datastore.ScanningCondition{
				PrimaryKeyNs:  entities.IdNsEp,
				PrimaryKeyCol: "req_id",
			},
		)

		var data []entities.Attempt
		if tx := tx.Find(&data); tx.Error != nil {
			ch <- &datastore.ScanningRecord[[]entities.Attempt]{Error: tx.Error}
			return
		}

		ch <- &datastore.ScanningRecord[[]entities.Attempt]{Data: data}

		if len(data) < query.Size {
			return
		}

		// refresh cursor for next round
		// by default datastore.ScanningOrderDesc will be used, so the last item must be use as cursor
		cursor = data[len(data)-1].ReqId
	}
}

func (sql *SqlAttempt) ListRequests(ctx context.Context, attempts map[string]*entities.Attempt) (map[string]*entities.Request, error) {
	returning := map[string]*entities.Request{}
	if len(attempts) == 0 {
		return returning, nil
	}

	refs := map[string]string{}
	conditions := []string{}
	values := map[string]any{}
	var i int
	for refId := range attempts {
		refs[attempts[refId].ReqId] = refId
		values[fmt.Sprintf("ep_id_%d", i)] = attempts[refId].EpId
		values[fmt.Sprintf("msg_id_%d", i)] = attempts[refId].MsgId
		values[fmt.Sprintf("id_%d", i)] = attempts[refId].ReqId
		conditions = append(conditions, fmt.Sprintf("(ep_id = @ep_id_%d AND msg_id = @msg_id_%d AND id = @id_%d)", i, i, i))
		i++
	}

	tableName := fmt.Sprintf(`"%s"`, entities.TableReq)
	statement := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s;",
		tableName,
		strings.Join(conditions, " OR "),
	)

	rows, err := sql.client.Raw(statement, values).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var req entities.Request
		if err := sql.client.ScanRows(rows, &req); err != nil {
			return nil, err
		}

		returning[refs[req.Id]] = &req
	}

	return returning, nil
}

func (sql *SqlAttempt) Update(ctx context.Context, updates map[string]*entities.AttemptState) map[string]error {
	ko := safe.Map[error]{}
	var wg conc.WaitGroup
	for id := range updates {
		reqId := id
		wg.Go(func() {
			tx := sql.client.Model(&entities.Attempt{}).Where("req_id = ?", reqId).Updates(updates[reqId])
			if tx.Error != nil {
				ko.Set(reqId, tx.Error)
			}
		})
	}
	wg.Wait()
	return ko.Data()
}
