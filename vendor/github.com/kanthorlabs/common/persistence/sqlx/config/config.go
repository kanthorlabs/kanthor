package config

import "github.com/kanthorlabs/common/validator"

var (
	TypePostgres = "postgres"
	TypeSqlite   = "file"
)

type Config struct {
	Uri            string     `json:"uri" yaml:"uri" mapstructure:"uri"`
	SkipDefaultTxn bool       `json:"skip_default_txn" yaml:"skip_default_txn" mapstructure:"skip_default_txn"`
	Connection     Connection `json:"connection" yaml:"connection" mapstructure:"connection"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.StringUri("SQLX.CONFIG.URI", conf.Uri),
		validator.StringStartsWithOneOf("SQLX.CONFIG.URI", conf.Uri, []string{TypePostgres, TypeSqlite}),
	)
	if err != nil {
		return err
	}

	if err := conf.Connection.Validate(); err != nil {
		return err
	}

	return nil
}

type Connection struct {
	MaxLifetime  int64 `json:"max_lifetime" yaml:"max_lifetime" mapstructure:"max_lifetime"`
	MaxIdletime  int64 `json:"max_idletime" yaml:"max_idletime" mapstructure:"max_idletime"`
	MaxIdleCount int   `json:"max_idle_count" yaml:"max_idle_count" mapstructure:"max_idle_count"`
	MaxOpenCount int   `json:"max_open_count" yaml:"max_open_count" mapstructure:"max_open_count"`
}

func (conf *Connection) Validate() error {
	return validator.Validate(
		validator.NumberInRange("SQLX.CONFIG.CONNECTION.MAX_LIFETIME", conf.MaxLifetime, 300000, 3600000),
		validator.NumberInRange("SQLX.CONFIG.CONNECTION.MAX_IDLETIME", conf.MaxIdletime, 60000, 3600000),
		validator.NumberGreaterThan("SQLX.CONFIG.CONNECTION.MAX_IDLE_COUNT", conf.MaxIdleCount, 1),
		validator.NumberGreaterThan("SQLX.CONFIG.CONNECTION.MAX_OPEN_COUNT", conf.MaxOpenCount, 1),
	)
}
