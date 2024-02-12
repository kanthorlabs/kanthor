package config

import (
	"slices"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/services"
	attempt "github.com/kanthorlabs/kanthor/services/attempt/config"
	dispatcher "github.com/kanthorlabs/kanthor/services/dispatcher/config"
	portal "github.com/kanthorlabs/kanthor/services/portal/config"
	recovery "github.com/kanthorlabs/kanthor/services/recovery/config"
	scheduler "github.com/kanthorlabs/kanthor/services/scheduler/config"
	sdk "github.com/kanthorlabs/kanthor/services/sdk/config"
	storage "github.com/kanthorlabs/kanthor/services/storage/config"
)

func Services(provider configuration.Provider, name string) (map[string]validator.Validator, error) {
	returning := map[string]validator.Validator{}
	if name == services.SDK || name == services.ALL {
		s, err := sdk.New(provider)
		if err != nil {
			return nil, err
		}
		returning[services.SDK] = s
	}
	if name == services.PORTAL || name == services.ALL {
		s, err := portal.New(provider)
		if err != nil {
			return nil, err
		}
		returning[services.PORTAL] = s
	}
	if name == services.SCHEDULER || name == services.ALL {
		s, err := scheduler.New(provider)
		if err != nil {
			return nil, err
		}
		returning[services.SCHEDULER] = s
	}
	if name == services.DISPATCHER || name == services.ALL {
		s, err := dispatcher.New(provider)
		if err != nil {
			return nil, err
		}
		returning[services.DISPATCHER] = s
	}
	if name == services.STORAGE || name == services.ALL {
		s, err := storage.New(provider)
		if err != nil {
			return nil, err
		}
		returning[services.STORAGE] = s
	}
	if slices.Contains(services.SERVICE_RECOVERY, name) || name == services.ALL {
		s, err := recovery.New(provider)
		if err != nil {
			return nil, err
		}
		returning[services.STORAGE] = s
	}
	if slices.Contains(services.SERVICE_ATTEMPT, name) || name == services.ALL {
		s, err := attempt.New(provider)
		if err != nil {
			return nil, err
		}
		returning[services.STORAGE] = s
	}

	return returning, nil
}
