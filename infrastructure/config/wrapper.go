package config

import (
	"github.com/kanthorlabs/common/configuration"
)

func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	if err := provider.Unmarshal(&conf); err != nil {
		return nil, err
	}
	if err := conf.Infrastructure.Validate(); err != nil {
		return nil, err
	}

	return &conf.Infrastructure, nil
}

type Wrapper struct {
	Infrastructure Config `json:"infrastructure" yaml:"infrastructure" mapstructure:"infrastructure"`
}
