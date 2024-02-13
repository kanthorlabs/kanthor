package db

import (
	"context"
	"fmt"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SqlEndpoint struct {
	client *gorm.DB
}

func (sql *SqlEndpoint) Create(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error) {
	if tx := sql.client.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpoint) Update(ctx context.Context, doc *entities.Endpoint) (*entities.Endpoint, error) {
	updates := []string{
		"name",
		"method",
		"uri",
		"updated_at",
	}
	tx := sql.client.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		// When update with struct, GORM will only update non-zero fields,
		// you might want to use map to update attributes or use Select to specify fields to update
		Select(updates).
		Updates(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpoint) Delete(ctx context.Context, doc *entities.Endpoint) error {
	tx := sql.client.WithContext(ctx).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Delete(doc)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlEndpoint) List(ctx context.Context, wsId, appId string, query *database.PagingQuery) ([]entities.Endpoint, error) {
	doc := &entities.Endpoint{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseAppId(appId, doc.TableName()),
			UseWsId(wsId, entities.TableApp),
		).
		Order(clause.OrderByColumn{Column: clause.Column{Name: fmt.Sprintf(`"%s"."created_at"`, doc.TableName())}, Desc: true})

	tx = query.Sqlx(tx, "id", []string{doc.ColName("name"), doc.ColName("uri")})

	var docs []entities.Endpoint
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlEndpoint) Count(ctx context.Context, wsId, appId string, query *database.PagingQuery) (int64, error) {
	doc := &entities.Endpoint{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseAppId(appId, doc.TableName()),
			UseWsId(wsId, entities.TableApp),
		)

	tx = query.SqlxCount(tx, "id", []string{doc.ColName("name"), doc.ColName("uri")})

	var count int64
	return count, tx.Count(&count).Error
}

func (sql *SqlEndpoint) Get(ctx context.Context, wsId string, id string) (*entities.Endpoint, error) {
	doc := &entities.Endpoint{}
	doc.Id = id

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseApp(doc.TableName()),
			UseWsId(wsId, entities.TableApp),
		).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
