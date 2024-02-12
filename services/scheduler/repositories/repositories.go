package repositories

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/database"
	"github.com/kanthorlabs/kanthor/services/scheduler/repositories/db"
)

func New(logger logging.Logger, dbclient database.Database) Repositories {
	return NewSql(logger, dbclient)
}

type Repositories interface {
	Database() db.Database
}
