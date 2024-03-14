package entities

import (
	"fmt"
	"testing"

	"github.com/kanthorlabs/common/cipher/encryption"
	"github.com/kanthorlabs/common/utils"
	"github.com/stretchr/testify/require"
)

func TestEndpoint_RotateSecretKey(t *testing.T) {
	encryptkey := utils.RandomString(32)

	t.Run("OK", func(t *testing.T) {
		entity := &Endpoint{}
		require.NoError(t, entity.RotateSecretKey(encryptkey, 64))
		require.NotEmpty(t, entity.SecretKey)

		// rotate
		secretKey := entity.SecretKey
		require.NoError(t, entity.RotateSecretKey(encryptkey, 64))
		require.NotEmpty(t, entity.SecretKey)
		require.NotEqual(t, secretKey, entity.SecretKey)
	})

	t.Run("KO - length error", func(t *testing.T) {
		entity := &Endpoint{}
		require.Error(t, entity.RotateSecretKey(encryptkey, 16))
		require.Empty(t, entity.SecretKey)
		require.Error(t, entity.RotateSecretKey(encryptkey, 65))
		require.Empty(t, entity.SecretKey)
	})

	t.Run("KO - generate new secret key error", func(t *testing.T) {
		entity := &Endpoint{}
		require.ErrorContains(t, entity.RotateSecretKey("", 64), "ENCRIPTION.ENCRYPT")
	})

	t.Run("KO - rotate secret key error", func(t *testing.T) {
		entity := &Endpoint{}
		require.NoError(t, entity.RotateSecretKey(encryptkey, 64))
		require.NotEmpty(t, entity.SecretKey)
		require.ErrorContains(t, entity.RotateSecretKey("", 64), "ENCRIPTION.DECRYPT")
	})
}

func TestEndpoint_DescryptSecretKey(t *testing.T) {
	encryptkey := utils.RandomString(32)

	t.Run("OK", func(t *testing.T) {
		secret := fmt.Sprintf("1=%s", utils.RandomString(64))
		encryptedsecret, err := encryption.Encrypt(encryptkey, secret)
		require.NoError(t, err)

		entity := &Endpoint{SecretKey: encryptedsecret}

		version, key, err := entity.DescryptSecretKey(encryptkey)
		require.NoError(t, err)
		require.Equal(t, secret, fmt.Sprintf("%d=%s", version, key))
	})

	t.Run("KO - decrypt error", func(t *testing.T) {
		entity := &Endpoint{}

		_, _, err := entity.DescryptSecretKey("")
		require.ErrorContains(t, err, "ENCRIPTION.DECRYPT")
	})

	t.Run("KO - secret key format error", func(t *testing.T) {
		secret := utils.RandomString(64)
		encryptedsecret, err := encryption.Encrypt(encryptkey, secret)
		require.NoError(t, err)

		entity := &Endpoint{SecretKey: encryptedsecret}

		_, _, err = entity.DescryptSecretKey(encryptkey)
		require.ErrorContains(t, err, "invalid secret key format")
	})

	t.Run("KO - secret key version error", func(t *testing.T) {
		secret := fmt.Sprintf("xxxxxxxxxxxxxxxxxxxxxx=%s", utils.RandomString(64))
		encryptedsecret, err := encryption.Encrypt(encryptkey, secret)
		require.NoError(t, err)

		entity := &Endpoint{SecretKey: encryptedsecret}

		_, _, err = entity.DescryptSecretKey(encryptkey)
		require.ErrorContains(t, err, "invalid secret key version")
	})
}
