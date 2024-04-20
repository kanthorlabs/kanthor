package config

import (
	gateway "github.com/kanthorlabs/common/gateway/config"
	"github.com/kanthorlabs/common/validator"
)

var ServiceName = "sdk"

type Sdk struct {
	Authn   Authn          `json:"authn" yaml:"authn" mapstructure:"authn"`
	Gateway gateway.Config `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
}

func (c *Sdk) Validate() error {
	if err := c.Authn.Validate(); err != nil {
		return err
	}
	if err := c.Gateway.Validate(); err != nil {
		return err
	}

	return nil
}

type Authn struct {
	DefaultStrategy string `json:"default_strategy" yaml:"default_strategy" mapstructure:"default_strategy"`
}

func (c *Authn) Validate() error {
	return validator.Validate(validator.StringRequired("SDK.CONFIG.AUTHN.DEFAULT_STRATEGY", c.DefaultStrategy))
}
