package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type Encryptor interface {
	EncryptString(data string) (string, error)
	DecryptString(data string) (string, error)
}

func NewAes(key string) (Encryptor, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	return &Aes{block: block}, nil
}

type Aes struct {
	block cipher.Block
}

func (encryptor *Aes) EncryptString(data string) (string, error) {
	b := []byte(data)
	// ciphertext includes both nonce & encrypted
	ciphertext := make([]byte, aes.BlockSize+len(b))

	// fill random nonce
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.New("CIPHER.ENCRYPTOR.ENCRYPT_STRING.NONCE_GENERATE.ERROR")
	}

	stream := cipher.NewCFBEncrypter(encryptor.block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], b)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (encryptor *Aes) DecryptString(data string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.New("CIPHER.ENCRYPTOR.DECRYPT_STRING.DATA_INVALID.ERROR")
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("CIPHER.ENCRYPTOR.DECRYPT_STRING.DATA_LENGTH_INVALID.ERROR")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(encryptor.block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
