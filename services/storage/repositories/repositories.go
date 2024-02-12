package repositories

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/datastore"
	"github.com/kanthorlabs/kanthor/services/storage/repositories/ds"
)

func New(logger logging.Logger, dsclient datastore.Datastore) Repositories {
	return NewSql(logger, dsclient)
}

type Repositories interface {
	Datastore() ds.Datastore
}
