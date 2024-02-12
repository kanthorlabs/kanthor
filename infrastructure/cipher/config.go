package cipher

import "github.com/kanthorlabs/common/validator"

type Config struct {
	Secret string `json:"secret" yaml:"secret" mapstructure:"secret"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.StringLen("CIPHER.CONFIG.SECRET", conf.Secret, 32, 32),
	)
}
