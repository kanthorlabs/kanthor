package usecase

import (
	"sync"

	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/services/recovery/config"
	"github.com/kanthorlabs/kanthor/services/recovery/repositories"
)

type Recovery interface {
	Scanner() Scanner
}

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	repositories repositories.Repositories,
) Recovery {
	logger = logger.With("usecase", "recovery")

	return &recovery{
		conf:         conf,
		logger:       logger,
		infra:        infra,
		repositories: repositories,
	}
}

type recovery struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories

	scanner *scanner

	mu sync.Mutex
}

func (uc *recovery) Scanner() Scanner {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.scanner == nil {
		uc.scanner = &scanner{
			conf:         uc.conf,
			logger:       uc.logger,
			infra:        uc.infra,
			publisher:    uc.infra.Stream.Publisher("recovery.scanner"),
			repositories: uc.repositories,
		}
	}
	return uc.scanner
}
