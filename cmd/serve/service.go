package serve

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services"
	"github.com/kanthorlabs/kanthor/services/ioc"
)

func Service(provider configuration.Provider, name string) (patterns.Runnable, error) {
	if name == services.PORTAL {
		return ioc.Portal(provider)
	}
	if name == services.SDK {
		return ioc.Sdk(provider)
	}
	if name == services.SCHEDULER {
		return ioc.Scheduler(provider)
	}
	if name == services.DISPATCHER {
		return ioc.Dispatcher(provider)
	}
	if name == services.STORAGE {
		return ioc.Storage(provider)
	}
	if name == services.RECOVERY_CRONJOB {
		return ioc.RecoveryCronjob(provider)
	}
	if name == services.RECOVERY_CONSUMER {
		return ioc.RecoveryConsumer(provider)
	}
	if name == services.ATTEMPT_CRONJOB {
		return ioc.AttemptCronjob(provider)
	}
	if name == services.ATTEMPT_CONSUMER {
		return ioc.AttemptConsumer(provider)
	}
	if name == services.ATTEMPT_TRIGGER {
		return ioc.AttemptTrigger(provider)
	}
	if name == services.ATTEMPT_SELECTOR {
		return ioc.AttemptSelector(provider)
	}
	if name == services.ATTEMPT_ENDEAVOR {
		return ioc.AttemptEndeavor(provider)
	}

	return nil, fmt.Errorf("serve.service: unknown service [%s]", name)
}

func Services(provider configuration.Provider, names []string) ([]patterns.Runnable, error) {
	instances := []patterns.Runnable{}

	for _, name := range services.SERVICES {
		if slices.Contains(names, services.ALL) || slices.Contains(names, name) {
			instance, err := Service(provider, name)
			if err != nil {
				return nil, err
			}

			instances = append(instances, instance)
		}

	}

	return instances, nil
}

type Stoppable interface {
	Stop(ctx context.Context) error
}

func Stop(instances ...Stoppable) error {
	// wait a little to stop our service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errc := make(chan error, 1)
	defer close(errc)
	go func() {
		var returning error
		for _, instance := range instances {
			if err := instance.Stop(ctx); err != nil {
				returning = errors.Join(returning, err)
			}
		}

		errc <- returning
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
