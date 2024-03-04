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

func multiple(provider configuration.Provider, names []string) error {
	logger, err := logging.New(provider)
	if err != nil {
		return err
	}

	services, err := Services(provider, names)
	if err != nil {
		return err
	}

	defer func() {
		stopctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		p := pool.New().WithContext(stopctx)
		for k := range services {
			name := k
			service := services[name]

			p.Go(func(subctx context.Context) error {
				return service.Stop(subctx)
			})
		}

		if err := p.Wait(); err != nil {
			logger.Error(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for k := range services {
		name := k
		service := services[name]

		if err = service.Start(ctx); err != nil {
			return err
		}
		go func() {
			if err = service.Run(ctx); err != nil {
				logger.Error(err)
			}
		}()
	}

	// listen for the interrupt signal.
	<-ctx.Done()
	logger.Infow("SERVER.RUN.INTERRUPT", "signal", fmt.Sprintf("%v", ctx))
	return nil
}
