package cipher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type Signature interface {
	SignString(data, key string) string
	VerifyString(data, key, compare string) error
}

func NewHmac() (Signature, error) {
	return &Hmac{}, nil
}

type Hmac struct {
}

func (signature *Hmac) SignString(data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))

	return hex.EncodeToString(mac.Sum(nil))
}

func (signature *Hmac) VerifyString(data, key, compare string) error {
	signed := signature.SignString(data, key)

	if hmac.Equal([]byte(signed), []byte(compare)) {
		return nil
	}

	return errors.New("CIPHER.SIGNATURE.VERIFY.NOT_MATCH.ERROR")
}
