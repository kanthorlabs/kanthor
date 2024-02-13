package ds

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/datastore"
)

func New(logger logging.Logger, ds datastore.Datastore) Datastore {
	return NewSql(logger, ds)
}

type Datastore interface {
	Message() Message
	Request() Request
	Response() Response
	Attempt() Attempt
}
