package db

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/database"
)

func New(logger logging.Logger, db database.Database) Database {
	return NewSql(logger, db)
}

type Database interface {
	Workspace() Workspace
	WorkspaceCredentials() WorkspaceCredentials
	Application() Application
	Endpoint() Endpoint
	EndpointRule() EndpointRule
}
