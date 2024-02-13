package db

import (
	"context"
	"fmt"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlEndpoint struct {
	client *gorm.DB
}

func (sql *SqlEndpoint) CreateBatch(ctx context.Context, docs []entities.Endpoint) ([]string, error) {
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

func (sql *SqlEndpoint) Count(ctx context.Context, wsId string, query *database.PagingQuery) (int64, error) {
	doc := &entities.Endpoint{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseApp(doc.TableName()),
			UseWsId(wsId, entities.TableApp),
		)

	tx = query.SqlxCount(tx, "id", []string{doc.ColName("name"), doc.ColName("uri")})

	var count int64
	return count, tx.Count(&count).Error
}

func (sql *SqlEndpoint) Get(ctx context.Context, wsId, id string) (*entities.Endpoint, error) {
	doc := &entities.Endpoint{}
	doc.Id = id

	tx := sql.client.WithContext(ctx).Model(&doc).
		Scopes(UseApp(doc.TableName()), UseWsId(wsId, entities.TableApp)).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
