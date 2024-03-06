package run

import (
	"fmt"
	"slices"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/kanthor/services/ioc"
	portalconfig "github.com/kanthorlabs/kanthor/services/portal/config"
)

var (
	ALL      = "all"
	SERVICES = []string{
		portalconfig.ServiceName,
	}
)

func Service(provider configuration.Provider, name string) (patterns.Runnable, error) {
	if name == portalconfig.ServiceName {
		return ioc.Portal(provider)
	}

	return nil, fmt.Errorf("SERVER.RUN.UNKNOWN_SERVICE.ERROR: [%s]", name)
}

func Services(provider configuration.Provider, names []string) (map[string]patterns.Runnable, error) {
	instances := map[string]patterns.Runnable{}

	for _, name := range SERVICES {
		if slices.Contains(names, ALL) || slices.Contains(names, name) {
			instance, err := Service(provider, name)
			if err != nil {
				return nil, err
			}

			instances[name] = instance
		}
	}

	return instances, nil
}
