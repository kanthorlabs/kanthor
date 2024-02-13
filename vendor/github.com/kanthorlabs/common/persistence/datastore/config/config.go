package config

import (
	"github.com/kanthorlabs/common/configuration"
	xqlxconf "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/validator"
)

var EngineSqlx = "sqlx"

func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	if err := provider.Unmarshal(&conf); err != nil {
		return nil, err
	}
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &conf.Datastore, nil
}

type Wrapper struct {
	Datastore Config `json:"datastore" yaml:"datastore" mapstructure:"datastore"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Datastore.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Engine string           `json:"engine" yaml:"engine" mapstructure:"engine"`
	Sqlx   *xqlxconf.Config `json:"sqlx" yaml:"sqlx" mapstructure:"sqlx"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.StringOneOf("DATASTORE.CONFIG.ENGINE", conf.Engine, []string{EngineSqlx}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineSqlx {
		if err := conf.Sqlx.Validate(); err != nil {
			return err
		}
	}

	return nil
}
