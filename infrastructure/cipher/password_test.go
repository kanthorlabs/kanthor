package cipher_test

import (
	"testing"

	"github.com/jaswdr/faker"
	"github.com/kanthorlabs/kanthor/infrastructure/cipher"
	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	password, err := cipher.NewBcrypt()
	assert.Nil(t, err)

	pass := faker.New().Internet().Password()
	hash, err := password.HashString(pass)
	assert.Nil(t, err)

	err = password.CompareString(pass, hash)
	assert.Nil(t, err)
}
