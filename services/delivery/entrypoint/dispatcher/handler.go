package dispatcher

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/streaming"
	stmentities "github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
	"github.com/kanthorlabs/kanthor/internal/transformation"
	"github.com/kanthorlabs/kanthor/services/delivery/usecase"
)

var ErrDispatcherForward = errors.New("DELIVERY.DISPATCHER.FORWARD.ERROR")

func handler(service *dispatcher) streaming.SubHandler {
	// If you return a map of errors, the message with the following refId will be retried
	return func(ctx context.Context, events map[string]*stmentities.Event) map[string]error {
		// the dispatcher has configured timeout for each request it makes
		// so don't need to set global timeout here

		in := &usecase.DispatcherForwardIn{
			Requests: make(map[string]*entities.Request),
		}

		for refId, event := range events {
			request, err := transformation.EventToRequest(event)
			if err != nil {
				service.logger.Errorw(ErrDispatcherForward.Error(), "error", err.Error(), "event", event.String())
				// unable to parse request from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			in.Requests[refId] = request
		}

		// we alreay validated messages of request, don't need to validate again
		out, err := service.uc.Dispatcher().Fowrard(ctx, in)
		// un-retriable error, reject the whole message batch
		if err != nil {
			service.logger.Errorw(ErrDispatcherForward.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return map[string]error{}
		}

		// any item of the .Error map will be retried
		return out.Error
	}
}
