package db

import (
	"context"
	"fmt"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
}

func (sql *SqlApplication) CreateBatch(ctx context.Context, docs []entities.Application) ([]string, error) {
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

func (sql *SqlApplication) Count(ctx context.Context, wsId string, query *database.PagingQuery) (int64, error) {
	doc := &entities.Application{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, doc.TableName()))

	tx = query.SqlxCount(tx, "id", []string{doc.ColName("name")})

	var count int64
	return count, tx.Count(&count).Error
}

func (sql *SqlApplication) Get(ctx context.Context, wsId, id string) (*entities.Application, error) {
	doc := &entities.Application{}
	doc.Id = id

	tx := sql.client.WithContext(ctx).Model(&doc).
		Scopes(UseWsId(wsId, doc.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
