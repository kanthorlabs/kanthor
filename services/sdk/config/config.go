package config

import (
	"github.com/kanthorlabs/common/configuration"
	logging "github.com/kanthorlabs/common/logging/config"
	infrastructure "github.com/kanthorlabs/kanthor/infrastructure/config"
)

func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	if err := provider.Unmarshal(&conf); err != nil {
		return nil, err
	}
	if err := conf.Logger.Validate(); err != nil {
		return nil, err
	}
	if err := conf.Infrastructure.Validate(); err != nil {
		return nil, err
	}

	return &conf, nil
}

type Config struct {
	Logger         logging.Config        `json:"logger" yaml:"logger" mapstructure:"logger"`
	Infrastructure infrastructure.Config `json:"infrastructure" yaml:"infrastructure" mapstructure:"infrastructure"`
	Sdk            Sdk                   `json:"sdk" yaml:"sdk" mapstructure:"sdk"`
}
