package repositories

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/database"
	"github.com/kanthorlabs/common/persistence/datastore"
	"github.com/kanthorlabs/kanthor/services/attempt/repositories/db"
	"github.com/kanthorlabs/kanthor/services/attempt/repositories/ds"
)

func New(logger logging.Logger, dbclient database.Database, dsclient datastore.Datastore) Repositories {
	return NewSql(logger, dbclient, dsclient)
}

type Repositories interface {
	Database() db.Database
	Datastore() ds.Datastore
}
