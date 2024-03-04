package config

import (
	cache "github.com/kanthorlabs/common/cache/config"
	circuitbreaker "github.com/kanthorlabs/common/circuitbreaker/config"
	distributedlockmanager "github.com/kanthorlabs/common/distributedlockmanager/config"
	idempotency "github.com/kanthorlabs/common/idempotency/config"
	database "github.com/kanthorlabs/common/persistence/database/config"
	datastore "github.com/kanthorlabs/common/persistence/datastore/config"
	streaming "github.com/kanthorlabs/common/streaming/config"
)

type Config struct {
	Database       database.Config               `json:"database" yaml:"database" mapstructure:"database"`
	Datastore      datastore.Config              `json:"datastore" yaml:"datastore" mapstructure:"datastore"`
	Stream         streaming.Config              `json:"stream" yaml:"stream" mapstructure:"stream"`
	Cache          cache.Config                  `json:"cache" yaml:"cache" mapstructure:"cache"`
	DLM            distributedlockmanager.Config `json:"distributed_lock_manager" yaml:"distributed_lock_manager" mapstructure:"distributed_lock_manager"`
	Idempotency    idempotency.Config            `json:"idempotency" yaml:"idempotency" mapstructure:"idempotency"`
	CircuitBreaker circuitbreaker.Config         `json:"circuit_breaker" yaml:"circuit_breaker" mapstructure:"circuit_breaker"`
}

func (conf *Config) Validate() error {
	if err := conf.Database.Validate(); err != nil {
		return err
	}
	if err := conf.Datastore.Validate(); err != nil {
		return err
	}
	if err := conf.Stream.Validate(); err != nil {
		return err
	}
	if err := conf.Cache.Validate(); err != nil {
		return err
	}
	if err := conf.DLM.Validate(); err != nil {
		return err
	}
	if err := conf.Idempotency.Validate(); err != nil {
		return err
	}
	if err := conf.CircuitBreaker.Validate(); err != nil {
		return err
	}
	return nil
}
