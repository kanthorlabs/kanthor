package repositories

import (
	"errors"

	"github.com/kanthorlabs/common/persistence/datastore"
	dsconfig "github.com/kanthorlabs/common/persistence/datastore/config"
)

func New(ds datastore.Datastore) (Repositories, error) {
	if ds.Engine() == dsconfig.EngineSqlx {
		return &sqlrepos{ds: ds}, nil
	}

	return nil, errors.New("DATATORE.ENGINE.UNKNOWN.ERROR")
}

type Repositories interface {
	Message() Message
	Request() Request
	Response() Response
}
