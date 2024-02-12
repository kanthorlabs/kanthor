package repositories

import (
	"sync"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/database"
	"github.com/kanthorlabs/kanthor/datastore"
	"github.com/kanthorlabs/kanthor/services/portal/repositories/db"
	"github.com/kanthorlabs/kanthor/services/portal/repositories/ds"
)

func NewSql(logger logging.Logger, dbclient database.Database, dsclient datastore.Datastore) Repositories {
	logger = logger.With("repositories", "sql")
	return &sql{logger: logger, db: db.NewSql(logger, dbclient), ds: ds.NewSql(logger, dsclient)}
}

type sql struct {
	logger logging.Logger
	db     db.Database
	ds     ds.Datastore
	mu     sync.Mutex
}

func (repo *sql) Database() db.Database {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.db
}

func (repo *sql) Datastore() ds.Datastore {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.ds
}
