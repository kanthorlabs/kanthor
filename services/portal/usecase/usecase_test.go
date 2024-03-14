package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/testify"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setup(t *testing.T) (*portal, func()) {
	provider, cleanup := testify.Setup(t)

	conf, err := config.New(provider)
	require.NoError(t, err)

	logger, err := logging.NewNoop()
	require.NoError(t, err)

	infra, err := infrastructure.New(&conf.Infrastructure, logger)
	require.NoError(t, err)

	uc := &portal{
		conf:   conf,
		logger: logger,
		watch:  clock.New(),
		infra:  infra,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	require.NoError(t, uc.infra.Connect(ctx))

	uc.infra.Database().Client().(*gorm.DB).AutoMigrate(
		&entities.Workspace{},
		&entities.Application{},
		&entities.Endpoint{},
		&entities.Route{},
	)

	var terminate = func() {
		require.NoError(t, uc.infra.Disconnect(ctx))
		cleanup()
	}
	return uc, terminate
}

func assertcount(t *testing.T, uc *portal, table string, id string, expected int64) {
	var actual int64
	err := uc.infra.Database().Client().(*gorm.DB).
		Table(table).
		Where("id = ?", id).
		Count(&actual).Error
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
