package config

import (
	gateway "github.com/kanthorlabs/common/gateway/config"
)

var ServiceName = "sdk"

type Sdk struct {
	Gateway gateway.Config `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
}
