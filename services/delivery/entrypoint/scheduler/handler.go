package scheduler

import (
	"context"

	"github.com/kanthorlabs/common/streaming"
	stmentities "github.com/kanthorlabs/common/streaming/entities"
)

func handler(service *scheduler) streaming.SubHandler {
	return func(ctx context.Context, events map[string]*stmentities.Event) map[string]error {
		return map[string]error{}
	}
}
