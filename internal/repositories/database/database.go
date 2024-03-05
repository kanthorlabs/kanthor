package database

import (
	"sync"

	"github.com/kanthorlabs/common/persistence/database"
	"gorm.io/gorm"
)

func New(db database.Database) Database {
	return &sqlxrepo{client: db.Client().(*gorm.DB)}
}

type Database interface {
	Workspace() Workspace
}

type sqlxrepo struct {
	client *gorm.DB

	workspace *sqlxws
	mu        sync.Mutex
}

func (repo *sqlxrepo) Workspace() Workspace {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.workspace == nil {
		repo.workspace = &sqlxws{client: repo.client}
	}

	return repo.workspace
}
