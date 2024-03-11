package usecase

import (
	"sync"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
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
}

type sdk struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure

	mu sync.Mutex
}
