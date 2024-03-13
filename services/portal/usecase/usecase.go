package usecase

import (
	"sync"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"gorm.io/gorm"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra infrastructure.Infrastructure,
	watch clock.Clock,
) (Portal, error) {
	uc := &portal{
		conf:   conf,
		logger: logger,
		watch:  watch,
		infra:  infra,
	}
	return uc, nil
}

type Portal interface {
	Workspace() Workspace
	Credentials() Credentials
}

type portal struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure

	workspace   *workspace
	credentials *credentials

	mu sync.Mutex
}

func (uc *portal) Workspace() Workspace {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.workspace == nil {
		uc.workspace = &workspace{
			conf:   uc.conf,
			logger: uc.logger,
			watch:  uc.watch,
			infra:  uc.infra,
			orm:    uc.infra.Database().Client().(*gorm.DB),
		}
	}

	return uc.workspace
}

func (uc *portal) Credentials() Credentials {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.credentials == nil {
		uc.credentials = &credentials{
			conf:   uc.conf,
			logger: uc.logger,
			watch:  uc.watch,
			infra:  uc.infra,
			orm:    uc.infra.Database().Client().(*gorm.DB),
		}
	}

	return uc.credentials
}
