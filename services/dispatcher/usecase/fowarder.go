package usecase

import (
	"context"

	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/infrastructure/streaming"
	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/services/dispatcher/config"
)

type Forwarder interface {
	Send(ctx context.Context, in *ForwarderSendIn) (*ForwarderSendOut, error)
}

type forwarder struct {
	conf      *config.Config
	logger    logging.Logger
	infra     *infrastructure.Infrastructure
	publisher streaming.Publisher
}
