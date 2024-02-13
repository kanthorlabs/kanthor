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

	endpoint *SqlEndpoint

	mu sync.Mutex
}

func (repo *sql) Endpoint() Endpoint {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.endpoint == nil {
		repo.endpoint = &SqlEndpoint{client: repo.db.Client().(*gorm.DB)}
	}

	return repo.endpoint
}
