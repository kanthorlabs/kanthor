package repositories

import (
	"github.com/kanthorlabs/kanthor/database"
	"github.com/kanthorlabs/kanthor/datastore"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/services/portal/repositories/db"
	"github.com/kanthorlabs/kanthor/services/portal/repositories/ds"
)

func New(logger logging.Logger, dbclient database.Database, dsclient datastore.Datastore) Repositories {
	return NewSql(logger, dbclient, dsclient)
}

type Repositories interface {
	Database() db.Database
	Datastore() ds.Datastore
}
