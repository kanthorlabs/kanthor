package db

import (
	"sync"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/database"
	"gorm.io/gorm"
)

func NewSql(logger logging.Logger, db database.Database) Database {
	logger = logger.With("repositories", "db.sql")
	return &sql{logger: logger, db: db}
}

type sql struct {
	logger logging.Logger
	db     database.Database

	workspace   *SqlWorkspace
	application *SqlApplication

	mu sync.Mutex
}

func (repo *sql) Workspace() Workspace {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.workspace == nil {
		repo.workspace = &SqlWorkspace{client: repo.db.Client().(*gorm.DB)}
	}

	return repo.workspace
}

func (repo *sql) Application() Application {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.application == nil {
		repo.application = &SqlApplication{client: repo.db.Client().(*gorm.DB)}
	}

	return repo.application
}
