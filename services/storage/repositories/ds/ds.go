package ds

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/datastore"
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
