package usecase

import (
	"sync"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
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
	uc := &delivery{
		conf:   conf,
		logger: logger,
		watch:  watch,
		infra:  infra,
	}

	return uc, nil
}

type Delivery interface {
	Scheduler() Scheduler
}

type delivery struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure

	scheduler *scheduler

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
