package db

import (
	"context"
	"fmt"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlEndpointRule struct {
	client *gorm.DB
}

func (sql *SqlEndpointRule) Create(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error) {
	if tx := sql.client.WithContext(ctx).Create(doc); tx.Error != nil {
		return nil, tx.Error
	}
	return doc, nil
}

func (sql *SqlEndpointRule) Update(ctx context.Context, doc *entities.EndpointRule) (*entities.EndpointRule, error) {
	updates := []string{
		"name",
		"priority",
		"exclusionary",
		"condition_source",
		"condition_expression",
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

func (sql *SqlEndpointRule) Delete(ctx context.Context, doc *entities.EndpointRule) error {
	tx := sql.client.WithContext(ctx).Model(doc).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		Delete(doc)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlEndpointRule) List(ctx context.Context, wsId, appId, epId string, query *database.PagingQuery) ([]entities.EndpointRule, error) {
	doc := &entities.EndpointRule{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseEpId(epId, doc.TableName()),
			UseAppId(appId, entities.TableEp),
			UseWsId(wsId, entities.TableApp),
		)

	tx = tx.Order(fmt.Sprintf(`"%s"."exclusionary" DESC, "%s"."priority" DESC, "%s"."id" DESC`, doc.TableName(), doc.TableName(), doc.TableName()))

	tx = query.Sqlx(tx, "id", []string{doc.ColName("name"), doc.ColName("condition_source")})

	var docs []entities.EndpointRule
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlEndpointRule) Count(ctx context.Context, wsId, appId, epId string, query *database.PagingQuery) (int64, error) {
	doc := &entities.EndpointRule{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseEpId(epId, doc.TableName()),
			UseAppId(appId, entities.TableEp),
			UseWsId(wsId, entities.TableApp),
		)

	tx = query.SqlxCount(tx, "id", []string{doc.ColName("name"), doc.ColName("condition_source")})

	var count int64
	return count, tx.Count(&count).Error
}

func (sql *SqlEndpointRule) Get(ctx context.Context, wsId string, id string) (*entities.EndpointRule, error) {
	doc := &entities.EndpointRule{}
	doc.Id = id

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(
			UseEp(doc.TableName()),
			UseApp(entities.TableEp),
			UseWsId(wsId, entities.TableApp),
		).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
