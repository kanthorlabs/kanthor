package db

import (
	"context"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlEndpointRule struct {
	client *gorm.DB
}

func (sql *SqlEndpointRule) CreateBatch(ctx context.Context, docs []entities.EndpointRule) ([]string, error) {
	ids := []string{}
	if len(docs) == 0 {
		return ids, nil
	}

	for _, doc := range docs {
		ids = append(ids, doc.Id)
	}

	if tx := sql.client.WithContext(ctx).Create(docs); tx.Error != nil {
		return nil, tx.Error
	}
	return ids, nil
}

func (sql *SqlEndpointRule) Count(ctx context.Context, wsId string, query *database.PagingQuery) (int64, error) {
	doc := &entities.EndpointRule{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseEp(doc.TableName()),
			UseApp(entities.TableEp),
			UseWsId(wsId, entities.TableApp),
		)

	tx = query.SqlxCount(tx, "id", []string{doc.ColName("name"), doc.ColName("condition_source")})

	var count int64
	return count, tx.Count(&count).Error
}
