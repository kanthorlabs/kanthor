package usecase

import (
	"context"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/internal/datastore/repositories"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"gorm.io/gorm"
)

type Message interface {
	Create(ctx context.Context, in *MessageCreateIn) (*MessageCreateOut, error)
	Get(ctx context.Context, in *MessageGetIn) (*MessageGetOut, error)
}

type message struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	orm    *gorm.DB
	repos  repositories.Repositories
}
