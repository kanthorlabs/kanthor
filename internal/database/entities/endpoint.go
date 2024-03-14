package entities

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kanthorlabs/common/cipher/encryption"
	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
)

type Endpoint struct {
	Auditable

	// SecretKey format: version=secret_key
	SecretKey string

	AppId  string
	Name   string
	Method string
	Uri    string
}

func (entity *Endpoint) SetId() {
	entity.Id = idx.New(IdNsEp)
}

func (entity *Endpoint) TableName() string {
	return TableEp
}

func (entity *Endpoint) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("id", entity.Id, IdNsEp),
		validator.StringLen("secret_key", entity.SecretKey, 16, 32),
		validator.StringStartsWith("app_id", entity.AppId, IdNsApp),
		validator.StringRequired("name", entity.Name),
		validator.StringOneOf("method", entity.Method, []string{http.MethodPost, http.MethodPut}),
		validator.StringUri("uri", entity.Uri),
	)
}

func (entity *Endpoint) PrimaryProp() string {
	return fmt.Sprintf("%s.id", TableEp)
}

func (entity *Endpoint) SearchProps() []string {
	return []string{
		fmt.Sprintf("%s.name", TableEp),
		fmt.Sprintf("%s.method", TableEp),
		fmt.Sprintf("%s.uri", TableEp),
	}
}

func (entity *Endpoint) RotateSecretKey(encryptkey string, length int) error {
	lenok := length >= 32 && length <= 64
	if !lenok {
		return fmt.Errorf("invalid length: %d", length)
	}

	// if the secret key is empty, we should generate a new one with version 1
	if entity.SecretKey == "" {
		secret := fmt.Sprintf("1=%s", utils.RandomString(length))
		encryptedsecret, err := encryption.Encrypt(encryptkey, secret)
		if err != nil {
			return err
		}

		entity.SecretKey = encryptedsecret
		return nil
	}

	version, _, err := entity.DescryptSecretKey(encryptkey)
	if err != nil {
		return err
	}

	newsecret := fmt.Sprintf("%d=%s", version+1, utils.RandomString(length))
	encryptedsecret, err := encryption.Encrypt(encryptkey, newsecret)
	if err != nil {
		return err
	}

	entity.SecretKey = encryptedsecret
	return nil
}

func (entity *Endpoint) DescryptSecretKey(encryptkey string) (int, string, error) {
	secret, err := encryption.Decrypt(encryptkey, entity.SecretKey)
	if err != nil {
		return 0, "", err
	}

	parts := strings.Split(secret, "=")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid secret key format: %s", secret)
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid secret key version: %s", parts[0])
	}

	return version, parts[1], nil
}
