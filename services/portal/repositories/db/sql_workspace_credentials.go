package db

import (
	"context"
	"fmt"

	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SqlWorkspaceCredentials struct {
	client *gorm.DB
}

func (sql *SqlWorkspaceCredentials) Create(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error) {
	if tx := sql.client.Create(doc); tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}

func (sql *SqlWorkspaceCredentials) Update(ctx context.Context, doc *entities.WorkspaceCredentials) (*entities.WorkspaceCredentials, error) {
	updates := []string{
		"name",
		"expired_at",
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

func (sql *SqlWorkspaceCredentials) List(ctx context.Context, wsId string, query *database.PagingQuery) ([]entities.WorkspaceCredentials, error) {
	doc := &entities.WorkspaceCredentials{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, doc.TableName())).
		Order(clause.OrderByColumn{Column: clause.Column{Name: fmt.Sprintf(`"%s"."created_at"`, doc.TableName())}, Desc: true})

	tx = query.Sqlx(tx, "id", []string{doc.ColName("name")})

	var docs []entities.WorkspaceCredentials
	if tx = tx.Find(&docs); tx.Error != nil {
		return nil, tx.Error
	}

	return docs, nil
}

func (sql *SqlWorkspaceCredentials) Count(ctx context.Context, wsId string, query *database.PagingQuery) (int64, error) {
	doc := &entities.WorkspaceCredentials{}

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, entities.TableWsc))

	tx = query.SqlxCount(tx, "id", []string{doc.ColName("name")})

	var count int64
	return count, tx.Count(&count).Error
}

func (sql *SqlWorkspaceCredentials) Get(ctx context.Context, wsId, id string) (*entities.WorkspaceCredentials, error) {
	doc := &entities.WorkspaceCredentials{}
	doc.Id = id

	tx := sql.client.WithContext(ctx).Model(doc).
		Scopes(UseWsId(wsId, doc.TableName())).
		Where(fmt.Sprintf(`"%s"."id" = ?`, doc.TableName()), doc.Id).
		First(doc)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}
