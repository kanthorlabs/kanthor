package db

import (
	"sync"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/database"
	"gorm.io/gorm"
)

func NewSql(logger logging.Logger, db database.Database) Database {
	logger = logger.With("repositories", "sql")
	return &sql{logger: logger, db: db}
}

type sql struct {
	logger logging.Logger
	db     database.Database

	workspace            *SqlWorkspace
	workspaceCredentials *SqlWorkspaceCredentials
	application          *SqlApplication
	endpoint             *SqlEndpoint
	endpointRule         *SqlEndpointRule

	mu sync.Mutex
}

func (repo *sql) Workspace() Workspace {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.workspace == nil {
		repo.workspace = &SqlWorkspace{client: repo.db.Client().(*gorm.DB)}
	}

	return repo.workspace
}

func (repo *sql) WorkspaceCredentials() WorkspaceCredentials {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.workspaceCredentials == nil {
		repo.workspaceCredentials = &SqlWorkspaceCredentials{client: repo.db.Client().(*gorm.DB)}
	}

	return repo.workspaceCredentials
}

func (repo *sql) Application() Application {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.application == nil {
		repo.application = &SqlApplication{client: repo.db.Client().(*gorm.DB)}
	}

	return repo.application
}

func (repo *sql) Endpoint() Endpoint {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.endpoint == nil {
		repo.endpoint = &SqlEndpoint{client: repo.db.Client().(*gorm.DB)}
	}

	return repo.endpoint
}

func (repo *sql) EndpointRule() EndpointRule {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if repo.endpointRule == nil {
		repo.endpointRule = &SqlEndpointRule{client: repo.db.Client().(*gorm.DB)}
	}

	return repo.endpointRule
}
