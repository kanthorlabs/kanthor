package usecase

import (
	"sync"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/sender"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/delivery/config"
	"gorm.io/gorm"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra infrastructure.Infrastructure,
	watch clock.Clock,
) (Delivery, error) {
	send, err := sender.New(&conf.Dispatcher.Sender, logger)
	if err != nil {
		return nil, err
	}

	uc := &delivery{
		conf:   conf,
		logger: logger,
		watch:  watch,
		infra:  infra,
		send:   send,
	}

	return uc, nil
}

type Delivery interface {
	Scheduler() Scheduler
	Dispatcher() Dispatcher
}

type delivery struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	send   sender.Send

	scheduler  *scheduler
	dispatcher *dispatcher

	mu sync.Mutex
}

func (uc *delivery) Scheduler() Scheduler {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.scheduler == nil {
		uc.scheduler = &scheduler{
			conf:   uc.conf,
			logger: uc.logger,
			watch:  uc.watch,
			infra:  uc.infra,
			orm:    uc.infra.Database().Client().(*gorm.DB),
		}
	}

	return uc.scheduler
}

func (uc *delivery) Dispatcher() Dispatcher {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.dispatcher == nil {
		uc.dispatcher = &dispatcher{
			conf:   uc.conf,
			logger: uc.logger,
			watch:  uc.watch,
			infra:  uc.infra,
			send:   uc.send,
		}
	}

	return uc.dispatcher
}
