package cipher

import "golang.org/x/crypto/bcrypt"

type Password interface {
	HashString(pass string) (string, error)
	CompareString(pass, hash string) error
}

func NewBcrypt() (Password, error) {
	return &Bcrypt{}, nil
}

type Bcrypt struct {
}

func (password *Bcrypt) HashString(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

func (password *Bcrypt) CompareString(pass, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
