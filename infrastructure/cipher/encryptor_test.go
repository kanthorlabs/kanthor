package cipher_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/kanthorlabs/kanthor/infrastructure/cipher"
	"github.com/stretchr/testify/assert"
)

func TestEncryptor(t *testing.T) {
	key := strings.ReplaceAll(uuid.NewString(), "-", "")
	encryptor, err := cipher.NewAes(key)
	assert.Nil(t, err)

	data := faker.New().Lorem().Sentence(256)

	ciphertext, err := encryptor.EncryptString(data)
	assert.Nil(t, err)

	original, err := encryptor.DecryptString(ciphertext)
	assert.Nil(t, err)

	assert.Equal(t, data, original)
}
