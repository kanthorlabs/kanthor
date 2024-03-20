package config

import (
	sender "github.com/kanthorlabs/common/sender/config"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/constants"
)

var ServiceNameDispatcher = "dispatcher"

type Dispatcher struct {
	Topic  string        `json:"topic" yaml:"topic" mapstructure:"topic"`
	Sender sender.Config `json:"sender" yaml:"sender" mapstructure:"sender"`
}

func (conf *Dispatcher) Validate() error {
	err := validator.Validate(
		validator.StringStartsWithOneOf("DISPATCHER.CONFIG.TOPIC", conf.Topic, []string{
			constants.RequestTopic,
		}),
	)
	if err != nil {
		return err
	}

	if err := conf.Sender.Validate(); err != nil {
		return err
	}

	return nil
}
