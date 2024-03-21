package base

import (
	deliveryconfig "github.com/kanthorlabs/kanthor/services/delivery/config"
	portalconfig "github.com/kanthorlabs/kanthor/services/portal/config"
	sdkconfig "github.com/kanthorlabs/kanthor/services/sdk/config"
	storageconfig "github.com/kanthorlabs/kanthor/services/storage/config"
)

var ApiServiceNames = []string{
	portalconfig.ServiceName,
	sdkconfig.ServiceName,
}

var BackgroundServiceNames = []string{
	deliveryconfig.ServiceNameScheduler,
	deliveryconfig.ServiceNameDispatcher,
	storageconfig.ServiceName,
}

func ServiceNames() []string {
	var names []string
	names = append(names, ApiServiceNames...)
	names = append(names, BackgroundServiceNames...)
	return names
}
