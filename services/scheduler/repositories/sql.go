package repositories

import (
	"sync"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/database"
	"github.com/kanthorlabs/kanthor/services/scheduler/repositories/db"
)

func NewSql(logger logging.Logger, dbclient database.Database) Repositories {
	logger = logger.With("repositories", "sql")
	return &sql{logger: logger, db: db.NewSql(logger, dbclient)}
}

type sql struct {
	logger logging.Logger
	db     db.Database
	mu     sync.Mutex
}

func (repo *sql) Database() db.Database {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.db
}
