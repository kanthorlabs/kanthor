package usecase

import (
	"context"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/internal/datastore/repositories"
	"github.com/kanthorlabs/kanthor/services/storage/config"
)

type Message interface {
	Save(ctx context.Context, in *SaveMessageIn) (*SaveMessageOut, error)
}

type message struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	repos  repositories.Repositories
}
