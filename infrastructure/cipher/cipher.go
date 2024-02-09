package cipher

type Cipher struct {
	Password  Password
	Signature Signature
	Encryptor Encryptor
}

func New(conf *Config) (*Cipher, error) {
	password, err := NewBcrypt()
	if err != nil {
		return nil, err
	}

	signature, err := NewHmac()
	if err != nil {
		return nil, err
	}

	encryptor, err := NewAes(conf.Secret)
	if err != nil {
		return nil, err
	}

	return &Cipher{Password: password, Signature: signature, Encryptor: encryptor}, nil
}
