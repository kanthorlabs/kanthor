package database

import (
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/database/config"
	"github.com/kanthorlabs/kanthor/patterns"
)

func New(provider configuration.Provider) (Database, error) {
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

type Database interface {
	patterns.Connectable
	Client() any
}
