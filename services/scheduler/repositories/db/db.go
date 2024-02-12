package db

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/database"
)

func New(logger logging.Logger, db database.Database) Database {
	return NewSql(logger, db)
}

type Database interface {
	Application() Application
}
