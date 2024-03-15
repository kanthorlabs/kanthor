package config

import (
	gateway "github.com/kanthorlabs/common/gateway/config"
)

var ServiceName = "portal"

type Portal struct {
	Gateway gateway.Config `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
}

func (c *Portal) Validate() error {
	if err := c.Gateway.Validate(); err != nil {
		return err
	}

	return nil
}
