package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/logging"
	"github.com/sourcegraph/conc/pool"
)

func single(provider configuration.Provider, name string) error {
	logger, err := logging.New(provider)
	if err != nil {
		return err
	}

	service, err := Service(provider, name)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// start & stop
	if err = service.Start(ctx); err != nil {
		return err
	}
	defer func() {
		stopctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		p := pool.New().WithContext(stopctx)

		p.Go(func(subctx context.Context) error {
			return service.Stop(subctx)
		})

		if err := p.Wait(); err != nil {
			logger.Error(err)
		}
	}()

	go func() {
		if err = service.Run(ctx); err != nil {
			logger.Error(err)
		}
	}()

	// listen for the interrupt signal.
	<-ctx.Done()
	logger.Infow("SERVER.RUN.INTERRUPT", "signal", fmt.Sprintf("%v", ctx))
	return nil
}
