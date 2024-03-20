package usecase

import (
	"context"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/sender"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/delivery/config"
)

type Dispatcher interface {
	Fowrard(ctx context.Context, in *DispatcherForwardIn) (*DispatcherForwardOut, error)
}

type dispatcher struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	send   sender.Send
}
