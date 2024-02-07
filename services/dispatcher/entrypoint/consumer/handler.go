package consumer

import (
	"context"

	"github.com/kanthorlabs/kanthor/infrastructure/streaming"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"github.com/kanthorlabs/kanthor/internal/transformation"
	"github.com/kanthorlabs/kanthor/services/dispatcher/usecase"
)

func Handler(service *dispatcher) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(ctx context.Context, events map[string]*streaming.Event) map[string]error {
		in := &usecase.ForwarderSendIn{
			Concurrency: service.conf.Forwarder.Send.Concurrency,
			Requests:    make(map[string]*entities.Request),
		}

		for id, event := range events {
			request, err := transformation.EventToRequest(event)
			if err != nil {
				service.logger.Errorw("DISPATCHER.ENTRYPOINT.CONSUMER.HANDLER.EVENT_TRANSFORMATION.ERROR", "error", err.Error(), "event", event.String())
				// unable to parse request from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			if err := usecase.ValidateForwarderSendInRequest("request", request); err != nil {
				service.logger.Errorw("DISPATCHER.ENTRYPOINT.CONSUMER.HANDLER.REQUEST_VALIDATION.ERROR", "error", err.Error(), "event", event.String(), "request", request.String())
				// got malformed request, should ignore and not retry it
				continue
			}

			in.Requests[id] = request
		}

		// we alreay validated messages of request, don't need to validate again
		out, err := service.uc.Forwarder().Send(ctx, in)

		if err != nil {
			retruning := map[string]error{}
			// got un-coverable error, should retry all event
			for refId := range in.Requests {
				retruning[refId] = err
			}
			return retruning
		}

		return out.Error
	}
}
