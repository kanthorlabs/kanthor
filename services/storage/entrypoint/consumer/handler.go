package consumer

import (
	"context"
	"errors"
	"time"

	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/streaming"
	stmentities "github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/kanthor/internal/constants"
	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
	"github.com/kanthorlabs/kanthor/internal/transformation"
	"github.com/kanthorlabs/kanthor/services/storage/usecase"
)

var ErrStorageConsumerSave = errors.New("STORAGE.CONSUMER.SAVE.ERROR")

func handler(service *scheduler) streaming.SubHandler {
	// If you return a map of errors, the message with the following refId will be retried
	return func(ctx context.Context, events map[string]*stmentities.Event) map[string]error {
		requests := make(map[string]*entities.Request)
		responses := make(map[string]*entities.Response)

		for refId, event := range events {
			if project.IsTopic(event.Subject, constants.MessageTopic) {
				// every message must be put into datastore before publish to datastream
				// so we don't need to handle it in storage service
				continue
			}

			if project.IsTopic(event.Subject, constants.RequestTopic) {
				request, err := transformation.EventToRequest(event)
				if err != nil {
					service.logger.Errorw(ErrStorageConsumerSave.Error(), "error", err.Error(), "event", event.String())
					// unable to parse request from event is considered as un-retriable error
					// ignore the error, and we need to check it manually with log
					continue
				}

				requests[refId] = request
				continue
			}

			if project.IsTopic(event.Subject, constants.ResponseTopic) {
				response, err := transformation.EventToResponse(event)
				if err != nil {
					service.logger.Errorw(ErrStorageConsumerSave.Error(), "error", err.Error(), "event", event.String())
					// unable to parse response from event is considered as un-retriable error
					// ignore the error, and we need to check it manually with log
					continue
				}

				responses[refId] = response
				continue
			}

			// unknown topic, ignore the event
			service.logger.Warnw("STORAGE.CONSUMER.SAVE.UNKNOWN_TOPIC.ERROR", "event", event.String())
		}

		returning := safe.Map[error]{}

		// saving requests
		if len(requests) > 0 {
			timeout := time.Millisecond * time.Duration(service.conf.Storage.Message.Timeout)
			reqctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			in := &usecase.SaveRequestIn{Requests: requests}
			// we alreay validated requests items, don't need to validate again
			out, err := service.uc.Request().Save(reqctx, in)
			if err != nil {
				service.logger.Errorw(ErrStorageConsumerSave.Error(), "error", err.Error(), "requests", utils.Stringify(requests))
				// un-retriable error, reject the whole message batch
				return map[string]error{}
			}

			returning.Merge(out.Error)
		}

		// saving responses
		if len(responses) > 0 {
			timeout := time.Millisecond * time.Duration(service.conf.Storage.Message.Timeout)
			resctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			in := &usecase.SaveResponseIn{Responses: responses}
			// we alreay validated responses items, don't need to validate again
			out, err := service.uc.Response().Save(resctx, in)
			if err != nil {
				service.logger.Errorw(ErrStorageConsumerSave.Error(), "error", err.Error(), "responses", utils.Stringify(responses))
				// un-retriable error, reject the whole message batch
				return map[string]error{}
			}

			returning.Merge(out.Error)
		}

		return returning.Data()
	}
}
