package config

import (
	"fmt"

	"github.com/kanthorlabs/common/validator"
)

// Secrets contains all keys of cryptographic secrets used in the system.
// The first secret in the array is used for encryption data
// while all other keys are used to decrypt older data that were signed with.
type Secrets struct {
	Cipher []string `json:"cipher" yaml:"cipher" mapstructure:"cipher"`
}

func (conf *Secrets) Validate() error {
	return validator.Validate(
		validator.SliceRequired("INFRASTRUCTURE.CONFIG.SECRETS.CIPHER", conf.Cipher),
		validator.Slice(conf.Cipher, func(i int, item *string) error {
			key := fmt.Sprintf("INFRASTRUCTURE.CONFIG.SECRETS.CIPHER[%d]", i)
			return validator.StringLen(key, *item, 32, 32)()
		}),
	)
}
