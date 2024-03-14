package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/kanthor/internal/database/entities"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceCreate(t *testing.T) {
	uc, cleanup := setup(t)
	defer cleanup()

	t.Run("OK", func(t *testing.T) {
		in := &WorkspaceCreateIn{
			OwnerId: uuid.NewString(),
			Name:    testdata.Fake.App().Name(),
			Tier:    project.Tier(),
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		out, err := uc.Workspace().Create(ctx, in)

		require.NoError(t, err)
		require.NotEmpty(t, out.Id)
		require.Greater(t, out.CreatedAt, int64(0))
		require.Greater(t, out.UpdatedAt, int64(0))
		require.NotEmpty(t, out.Modifier)

		assertcount(t, uc, entities.TableWs, out.Id, 1)
	})
}
