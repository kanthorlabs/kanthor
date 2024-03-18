package usecase

import (
	"sync"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/internal/datastore/repositories"
	"github.com/kanthorlabs/kanthor/services/storage/config"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra infrastructure.Infrastructure,
	watch clock.Clock,
) (Storage, error) {
	repos, err := repositories.New(infra.Datastore())
	if err != nil {
		return nil, err
	}

	uc := &storage{
		conf:   conf,
		logger: logger,
		watch:  watch,
		infra:  infra,
		repos:  repos,
	}

	return uc, nil
}

type Storage interface {
	Message() Message
	Request() Request
	Response() Response
}

type storage struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	repos  repositories.Repositories

	message  *message
	request  *request
	response *response

	mu sync.Mutex
}

func (uc *storage) Message() Message {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.message == nil {
		uc.message = &message{
			conf:   uc.conf,
			logger: uc.logger,
			watch:  uc.watch,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}

	return uc.message
}

func (uc *storage) Request() Request {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.request == nil {
		uc.request = &request{
			conf:   uc.conf,
			logger: uc.logger,
			watch:  uc.watch,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}

	return uc.request
}

func (uc *storage) Response() Response {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.response == nil {
		uc.response = &response{
			conf:   uc.conf,
			logger: uc.logger,
			watch:  uc.watch,
			infra:  uc.infra,
			repos:  uc.repos,
		}
	}

	return uc.response
}
