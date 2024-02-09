package authenticator

import (
	"fmt"

	"github.com/kanthorlabs/kanthor/pkg/validator"
)

type Config struct {
	Engine string `json:"engine" yaml:"engine" mapstructure:"engine"`
	Ask    *Ask   `json:"ask" yaml:"ask" mapstructure:"ask"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf("AUTHENTICATOR.ENGINE", conf.Engine, []string{EngineAsk}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineAsk {
		return conf.Ask.Validate()
	}

	return nil
}

type Ask struct {
	Users []AskUser `json:"users" yaml:"users" mapstructure:"users"`
}

func (conf *Ask) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.SliceRequired("AUTHENTICATOR.ASK.USERS", conf.Users),
		validator.Slice(conf.Users, func(i int, item *AskUser) error {
			return item.Validate(fmt.Sprintf("AUTHENTICATOR.ASK.USERS[%d]", i))
		}),
	)
}

type AskUser struct {
	Username string `json:"username" yaml:"username" mapstructure:"username"`
	Password string `json:"password" yaml:"password" mapstructure:"password"`
}

func (conf *AskUser) Validate(prefix string) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired(prefix+".USERNAME", conf.Username),
		validator.StringRequired(prefix+".PASSWORD", conf.Password),
	)
}
