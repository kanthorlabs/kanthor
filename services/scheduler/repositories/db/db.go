package db

import (
	"github.com/kanthorlabs/kanthor/database"
	"github.com/kanthorlabs/kanthor/logging"
)

func New(logger logging.Logger, db database.Database) Database {
	return NewSql(logger, db)
}

type Database interface {
	Application() Application
}
