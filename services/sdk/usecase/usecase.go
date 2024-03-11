package usecase

import (
	"sync"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"gorm.io/gorm"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra infrastructure.Infrastructure,
	clock clock.Clock,
) (Sdk, error) {
	uc := &sdk{
		conf:   conf,
		logger: logger,
		watch:  clock,
		infra:  infra,
	}
	return uc, nil
}

type Sdk interface {
	Application() Application
}

type sdk struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure

	application *application

	mu sync.Mutex
}

func (uc *sdk) Application() Application {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.application == nil {
		uc.application = &application{
			conf:   uc.conf,
			logger: uc.logger,
			watch:  uc.watch,
			infra:  uc.infra,
			orm:    uc.infra.Database().Client().(*gorm.DB),
		}
	}

	return uc.application
}
