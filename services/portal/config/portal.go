package config

import (
	gateway "github.com/kanthorlabs/common/gateway/config"
)

type Portal struct {
	Gateway gateway.Config `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
}