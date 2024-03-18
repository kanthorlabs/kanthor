package config

import (
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/constants"
)

var ServiceName = "storage"

type Storage struct {
	Topic    string   `json:"topic" yaml:"topic" mapstructure:"topic"`
	Message  Message  `json:"message" yaml:"message" mapstructure:"message"`
	Request  Request  `json:"request" yaml:"request" mapstructure:"request"`
	Response Response `json:"response" yaml:"response" mapstructure:"response"`
}

func (conf *Storage) Validate() error {
	err := validator.Validate(
		validator.StringStartsWithOneOf("STORAGE.CONFIG.MESSAGE.TOPIC", conf.Topic, []string{
			constants.TopicCore,
			constants.MessageTopic,
			constants.RequestTopic,
			constants.ResponseTopic,
		}),
	)
	if err != nil {
		return err
	}

	if err := conf.Message.Validate(); err != nil {
		return err
	}
	if err := conf.Request.Validate(); err != nil {
		return err
	}
	if err := conf.Response.Validate(); err != nil {
		return err
	}

	return nil
}
