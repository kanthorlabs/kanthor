package db

import (
	"context"

	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"github.com/kanthorlabs/kanthor/internal/routing"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
}

func (sql *SqlApplication) Scan(ctx context.Context, query *datastore.ScanningQuery) chan *datastore.ScanningRecord[[]entities.Application] {
	ch := make(chan *datastore.ScanningRecord[[]entities.Application], 1)
	go sql.scan(ctx, query, ch)
	return ch
}

func (sql *SqlApplication) scan(ctx context.Context, query *datastore.ScanningQuery, ch chan *datastore.ScanningRecord[[]entities.Application]) {
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
			sql.client.Model(&entities.Application{}),
			&datastore.ScanningCondition{
				PrimaryKeyNs:  entities.IdNsEp,
				PrimaryKeyCol: "id",
			},
		)

		var data []entities.Application
		if tx := tx.Find(&data); tx.Error != nil {
			ch <- &datastore.ScanningRecord[[]entities.Application]{Error: tx.Error}
			return
		}

		ch <- &datastore.ScanningRecord[[]entities.Application]{Data: data}

		if len(data) < query.Size {
			return
		}

		// refresh cursor for next round
		// by default datastore.ScanningOrderDesc will be used, so the last item must be use as cursor
		cursor = data[len(data)-1].Id
	}
}

func (sql *SqlApplication) GetRoutes(ctx context.Context, ids []string) (map[string][]routing.Route, error) {
	returning := make(map[string][]routing.Route)

	endpoints, err := sql.getRouteEndpoints(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(endpoints) == 0 {
		return returning, nil
	}

	rules, err := sql.getRouteRules(ctx, endpoints)
	if err != nil {
		return nil, err
	}
	if len(rules) == 0 {
		return returning, nil
	}

	for i := range endpoints {
		if _, has := returning[endpoints[i].AppId]; has {
			returning[endpoints[i].AppId] = append(returning[endpoints[i].AppId], routing.Route{
				Endpoint: &endpoints[i],
				Rules:    rules[endpoints[i].Id],
			})
			continue
		}

		returning[endpoints[i].AppId] = []routing.Route{{
			Endpoint: &endpoints[i],
			Rules:    rules[endpoints[i].Id],
		}}
	}

	return returning, nil
}

func (sql *SqlApplication) getRouteEndpoints(ctx context.Context, appIds []string) ([]entities.Endpoint, error) {
	var endpoints []entities.Endpoint
	tx := sql.client.Model(&entities.Endpoint{}).Where("app_id IN ?", appIds).Find(&endpoints)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return endpoints, nil
}

func (sql *SqlApplication) getRouteRules(ctx context.Context, endpoints []entities.Endpoint) (map[string][]entities.EndpointRule, error) {
	returning := make(map[string][]entities.EndpointRule)
	var ids []string
	for i := range endpoints {
		returning[endpoints[i].Id] = make([]entities.EndpointRule, 0)
		ids = append(ids, endpoints[i].Id)
	}

	var rules []entities.EndpointRule
	tx := sql.client.
		Model(&entities.EndpointRule{}).
		Where("ep_id IN ?", ids).
		Order("ep_id DESC, exclusionary DESC, priority DESC").
		Find(&rules)
	if tx.Error != nil {
		return nil, tx.Error
	}

	for i := range rules {
		returning[rules[i].EpId] = append(returning[rules[i].EpId], rules[i])
	}
	return returning, nil
}
