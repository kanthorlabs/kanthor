package ds

import (
	"github.com/kanthorlabs/kanthor/datastore"
	"github.com/kanthorlabs/kanthor/logging"
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
