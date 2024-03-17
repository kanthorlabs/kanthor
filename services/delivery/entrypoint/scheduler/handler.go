package scheduler

import (
	"context"
	"errors"
	"time"

	"github.com/kanthorlabs/common/streaming"
	stmentities "github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/kanthor/internal/datastore/entities"
	"github.com/kanthorlabs/kanthor/internal/transformation"
	"github.com/kanthorlabs/kanthor/services/delivery/usecase"
)

var ErrSchedulerArrange = errors.New("DELIVERY.SCHEDULER.ARRANGE.ERROR")

func handler(service *scheduler) streaming.SubHandler {
	// If you return a map of errors, the message with the following refId will be retried
	return func(ctx context.Context, events map[string]*stmentities.Event) map[string]error {
		// It's important to make sure the subscriber handler timeout is set before processing the message
		timeout := time.Millisecond * time.Duration(service.conf.Scheduler.Request.Timeout)
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		in := &usecase.SchedulerArrangeIn{
			Messages: make(map[string]*entities.Message),
		}

		for refId, event := range events {
			message, err := transformation.EventToMessage(event)
			if err != nil {
				service.logger.Errorw(ErrSchedulerArrange.Error(), "error", err.Error(), "event", event.String())
				// unable to parse message from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			in.Messages[refId] = message
		}

		// we alreay validated messages of request, don't need to validate again
		out, err := service.uc.Scheduler().Arrange(ctx, in)
		// un-retriable error, reject the whole message batch
		if err != nil {
			service.logger.Errorw(ErrSchedulerArrange.Error(), "error", err.Error(), "in", utils.Stringify(in))
			return map[string]error{}
		}

		// any item of the .Error map will be retried
		return out.Error
	}
}
