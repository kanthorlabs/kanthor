package strategies

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/passport/utils"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/common/validator"
)

// NewAsk creates a new ask strategy instance what allows to authenticate users based on a username and password.
// The storage is based on a map of accounts in memory.
func NewAsk(conf *config.Ask, logger logging.Logger) (Strategy, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	accounts := make(map[string]*entities.Account)
	for i := range conf.Accounts {
		accounts[conf.Accounts[i].Username] = &conf.Accounts[i]
	}

	if len(accounts) != len(conf.Accounts) {
		return nil, errors.New("PASSPORT.STRATEGY.ASK.DUPLICATED_ACCOUNT.ERROR")
	}

	return &ask{conf: conf, logger: logger, accounts: accounts}, nil
}

type ask struct {
	conf   *config.Ask
	logger logging.Logger

	mu       sync.Mutex
	status   int
	accounts map[string]*entities.Account
}

func (instance *ask) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	instance.status = patterns.StatusConnected
	return nil
}

func (instance *ask) Readiness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return nil
}

func (instance *ask) Liveness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return nil
}

func (instance *ask) Disconnect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	instance.status = patterns.StatusDisconnected
	return nil
}

func (instance *ask) Register(ctx context.Context, acc entities.Account) error {
	return errors.New("PASSPORT.ASK.REGISTER.UNIMPLEMENT.ERROR")
}

func (instance *ask) Login(ctx context.Context, creds entities.Credentials) (*entities.Tokens, error) {
	return nil, errors.New("PASSPORT.ASK.LOGIN.UNIMPLEMENT.ERROR")
}

func (instance *ask) Logout(ctx context.Context, tokens entities.Tokens) error {
	return errors.New("PASSPORT.ASK.LOGOUT.UNIMPLEMENT.ERROR")
}

func (instance *ask) Verify(ctx context.Context, tokens entities.Tokens) (*entities.Account, error) {
	creds, err := utils.ParseBasicCredentials(tokens.Access)
	if err != nil {
		return nil, err
	}

	if err := creds.Validate(); err != nil {
		return nil, err
	}
	acc, has := instance.accounts[creds.Username]
	if !has {
		return nil, ErrLogin
	}

	if err := password.Compare(creds.Password, acc.PasswordHash); err != nil {
		return nil, ErrLogin
	}

	return acc.Censor(), nil
}

func (instance *ask) Management() Management {
	if instance.status != patterns.StatusConnected {
		panic(ErrNotConnected)
	}
	return &askmanagement{accounts: instance.accounts}
}

type askmanagement struct {
	accounts map[string]*entities.Account
}

func (instance *askmanagement) Deactivate(ctx context.Context, username string, ts int64) error {
	return errors.New("PASSPORT.ASK.MANAGEMENT.DEACTIVATE.UNIMPLEMENT.ERROR")
}

func (instance *askmanagement) List(ctx context.Context, usernames []string) ([]*entities.Account, error) {
	err := validator.Validate(
		validator.SliceRequired("usernames", usernames),
		validator.Slice(usernames, func(i int, item *string) error {
			key := fmt.Sprintf("usernames[%d]", i)
			return validator.StringRequired(key, *item)()
		}),
	)
	if err != nil {
		return nil, err
	}

	var accounts []*entities.Account
	for i := range usernames {
		if acc, has := instance.accounts[usernames[i]]; has {
			accounts = append(accounts, acc.Censor())
		}
	}

	return accounts, nil
}

func (instance *askmanagement) Update(ctx context.Context, acc entities.Account) error {
	return errors.New("PASSPORT.ASK.MANAGEMENT.UPDATE.UNIMPLEMENT.ERROR")
}
