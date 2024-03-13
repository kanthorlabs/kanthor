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
	Endpoint() Endpoint
	Route() Route
}

type sdk struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure

	application *application
	endpoint    *endpoint
	route       *route

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

func (uc *sdk) Endpoint() Endpoint {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.endpoint == nil {
		uc.endpoint = &endpoint{
			conf:   uc.conf,
			logger: uc.logger,
			watch:  uc.watch,
			infra:  uc.infra,
			orm:    uc.infra.Database().Client().(*gorm.DB),
		}
	}

	return uc.endpoint
}

func (uc *sdk) Route() Route {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.route == nil {
		uc.route = &route{
			conf:   uc.conf,
			logger: uc.logger,
			watch:  uc.watch,
			infra:  uc.infra,
			orm:    uc.infra.Database().Client().(*gorm.DB),
		}
	}

	return uc.route
}
