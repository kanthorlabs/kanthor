package repositories

import (
	"sync"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/services/storage/repositories/ds"
)

func NewSql(logger logging.Logger, dsclient datastore.Datastore) Repositories {
	logger = logger.With("repositories", "sql")
	return &sql{logger: logger, ds: ds.NewSql(logger, dsclient)}
}

type sql struct {
	logger logging.Logger
	ds     ds.Datastore
	mu     sync.Mutex
}

func (repo *sql) Datastore() ds.Datastore {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.ds
}
