package config

import "github.com/kanthorlabs/common/validator"

type Response struct {
	Timeout   int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	BatchSize int   `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
}

func (conf *Response) Validate() error {
	return validator.Validate(
		validator.NumberGreaterThanOrEqual("STORAGE.CONFIG.RESPONSE.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("STORAGE.CONFIG.RESPONSE.BATCH_SIZE", conf.BatchSize, 0),
	)
}
