package testify

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/project"
	natstest "github.com/nats-io/nats-server/v2/test"
	"github.com/stretchr/testify/require"
)

func Setup(t *testing.T) (configuration.Provider, func()) {
	_, fp, _, ok := runtime.Caller(0)
	require.True(t, ok)

	basepath := filepath.Dir(fp)

	opts := natstest.DefaultTestOptions
	opts.Port = 12224
	opts.JetStream = true
	opts.StoreDir = "/tmp/" + uuid.NewString()
	streaming := natstest.RunServer(&opts)

	provider, err := configuration.NewFile(project.Namespace(), []string{basepath})
	require.NoError(t, err)

	var cleanup = func() {
		streaming.Shutdown()
	}

	return provider, cleanup
}
