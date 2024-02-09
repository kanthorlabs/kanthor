package streaming

import (
	"context"

	"github.com/kanthorlabs/kanthor/logging"
	"github.com/kanthorlabs/kanthor/pkg/safe"
	"github.com/kanthorlabs/kanthor/pkg/utils"
	"github.com/kanthorlabs/kanthor/project"
	"github.com/kanthorlabs/kanthor/telemetry"
	natscore "github.com/nats-io/nats.go"
	"github.com/sourcegraph/conc/pool"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

type NatsPublisher struct {
	name   string
	conf   *Config
	logger logging.Logger

	nats *nats
}

func (publisher *NatsPublisher) Name() string {
	return publisher.name
}

func (publisher *NatsPublisher) Pub(ctx context.Context, events map[string]*Event) map[string]error {
	spanName := project.Topic("streaming.publisher.sub", publisher.name)
	spanner := ctx.Value(telemetry.CtxSpanner).(*telemetry.Spanner)

	datac := make(chan map[string]error, 1)
	defer close(datac)

	go func() {
		returning := safe.Map[error]{}
		p := pool.New().WithMaxGoroutines(publisher.conf.Publisher.RateLimit)
		for refId, e := range events {
			spanner.StartWithRefId(
				spanName, refId,
				attribute.String("streaming.publisher.engine", "nats"),
				attribute.String("event.id", refId),
				attribute.String("event.subject", e.Subject),
			)
			defer spanner.End(spanName)

			if err := e.Validate(); err != nil {
				publisher.logger.Errorw("STREAMING.PUBLISHER.NATS.EVENT_VALIDATION.ERROR", "subject", e.Subject, "event_id", e.Id, "event", e.String())
				returning.Set(refId, err)
				continue
			}

			// store the value to use in p.Go, otherwise we got the same value
			event := e
			msg := NatsMsgFromEvent(e.Subject, e)

			// telemetry tracing
			propgator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
			carrier := propagation.MapCarrier{}
			propgator.Inject(spanner.Contexts[refId], carrier)
			msg.Header.Add(HeaderTelemetryTrace, utils.Stringify(carrier))

			p.Go(func() {
				ack, err := publisher.nats.js.PublishMsg(msg, natscore.Context(ctx), natscore.MsgId(event.Id))
				if err != nil {
					publisher.logger.Errorw("STREAMING.PUBLISHER.NATS.EVENT_PUBLISH.ERROR", "subject", event.Subject, "event_id", event.Id)
					returning.Set(refId, err)
					return
				}

				if ack.Duplicate {
					publisher.logger.Errorw("STREAMING.PUBLISHER.NATS.EVENT_DUPLICATED.ERROR", "subject", event.Subject, "event_id", event.Id)
				}
			})
		}
		p.Wait()

		datac <- returning.Data()
	}()

	select {
	case data := <-datac:
		return data
	case <-ctx.Done():
		data := map[string]error{}
		for refId := range events {
			data[refId] = ctx.Err()
		}
		return data
	}
}
