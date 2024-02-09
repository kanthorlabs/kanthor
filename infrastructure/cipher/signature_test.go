package cipher_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/kanthorlabs/kanthor/infrastructure/cipher"
	"github.com/stretchr/testify/assert"
)

func TestSignature(t *testing.T) {
	signature, err := cipher.NewHmac()
	assert.Nil(t, err)

	data := faker.New().Lorem().Sentence(256)
	key := uuid.NewString()

	signed := signature.SignString(data, key)
	err = signature.VerifyString(data, key, signed)
	assert.Nil(t, err)
}
