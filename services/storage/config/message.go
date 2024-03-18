package config

import "github.com/kanthorlabs/common/validator"

type Message struct {
	Timeout   int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	BatchSize int   `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
}

func (conf *Message) Validate() error {
	return validator.Validate(
		validator.NumberGreaterThanOrEqual("STORAGE.CONFIG.MESSAGE.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("STORAGE.CONFIG.MESSAGE.BATCH_SIZE", conf.BatchSize, 0),
	)
}
