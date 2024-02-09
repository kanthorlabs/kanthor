package authenticator

import "context"

var EngineAsk = "ask"

func NewAsk(conf *Ask) (Verifier, error) {
	return &ask{conf: conf}, nil
}

type ask struct {
	conf *Ask
}

func (verifier *ask) Verify(ctx context.Context, request *Request) (*Account, error) {
	user, pass, err := ParseBasicCredentials(request.Credentials)
	if err != nil {
		return nil, err
	}

	for i := range verifier.conf.Users {
		if user != verifier.conf.Users[i].Username {
			continue
		}

		if pass != verifier.conf.Users[i].Password {
			return nil, ErrInvalidCredentials
		}

		return &Account{Sub: user, Name: user, Metadata: make(map[string]string, 0)}, nil
	}

	return nil, ErrInvalidCredentials
}
