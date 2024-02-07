package datastore

import (
	"github.com/kanthorlabs/kanthor/configuration"
	"github.com/kanthorlabs/kanthor/datastore/config"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/patterns"
)

func New(provider configuration.Provider) (Datastore, error) {
	conf, err := config.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}
	return NewSQL(conf, logger)
}

type Datastore interface {
	patterns.Connectable
	Client() any
}
