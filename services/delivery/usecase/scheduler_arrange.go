package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/cipher/encryption"
	"github.com/kanthorlabs/common/safe"
	stmentities "github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/conductor"
	"github.com/kanthorlabs/kanthor/internal/constants"
	dbentities "github.com/kanthorlabs/kanthor/internal/database/entities"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
	"github.com/kanthorlabs/kanthor/internal/transformation"
	"github.com/sourcegraph/conc"
)

var ErrSchedulerArrange = errors.New("DELIVERY.SCHEDULER.ARRANGE.ERROR")

func (uc *scheduler) Arrange(ctx context.Context, in *SchedulerArrangeIn) (*SchedulerArrangeOut, error) {
	ok := &safe.Slice[string]{}
	ko := &safe.Map[error]{}

	refs := &safe.Map[string]{}

	// arrange requests to events
	events := &safe.Map[*stmentities.Event]{}
	var wg conc.WaitGroup
	for id := range in.Messages {
		refId := id
		msg := in.Messages[refId]

		wg.Go(func() {
			reqs, err := uc.arrange(ctx, msg)
			if err != nil {
				uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", err, "message", utils.Stringify(msg))
				// if we got any error when arranging the message, reject the whole message because those error is unrecoverable
				return
			}

			var scheduleerr error
			schedulable := map[string]*stmentities.Event{}
			for _, req := range reqs {
				event, err := transformation.EventFromRequest(req, constants.RequestScheduleSubject)
				if err != nil {
					scheduleerr = errors.Join(scheduleerr, err)
					uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", err, "request", utils.Stringify(req))
					continue
				}

				schedulable[req.EpId] = event
			}

			// if we got any error when transforming the request to event, reject the whole message because those error is unrecoverable
			if scheduleerr != nil {
				for epId, event := range schedulable {
					uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", fmt.Sprintf("REJECTED by %v", scheduleerr), "event", event.String(), "ep_id", epId, "ref_id", refId)
				}
				return
			}

			for epId := range schedulable {
				// one message can produce multiple events, so if you set event key by refId, it will be overwritten
				// that why we need to create a direction of message -> refId -> epId -> event
				refs.Set(epId, refId)
				events.Set(epId, schedulable[epId])
			}
		})
	}
	wg.Wait()

	publisher, err := uc.infra.Streaming().Publisher(constants.RequestPublisher)
	if err != nil {
		uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", err.Error())
		return nil, ErrSchedulerArrange
	}

	errs := publisher.Pub(ctx, events.Data())
	if len(errs) > 0 {
		for epId, err := range errs {
			event, _ := events.Get(epId)
			uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", utils.Stringify(err), "event", event.String())

			// follow back of the direction message -> refId -> epId -> event
			// we report back to the streaming that we need to retry the message because one of request event got error
			refId, _ := refs.Get(epId)
			ko.Set(refId, err)
		}
	}

	return &SchedulerArrangeOut{Success: ok.Data(), Error: ko.Data()}, nil
}

func (uc *scheduler) arrange(ctx context.Context, msg *dsentities.Message) (map[string]*dsentities.Request, error) {
	destinations, err := uc.buildDestinationOfApp(ctx, msg.AppId)
	if err != nil {
		return nil, err
	}

	return conductor.Many(msg, destinations, uc.watch.Now())
}

func (uc *scheduler) buildDestinationOfApp(ctx context.Context, appId string) (map[string]*conductor.Destination, error) {
	endpoints, err := uc.getEndpointOfApp(ctx, appId)
	if err != nil {
		return nil, err
	}

	destinations := map[string]*conductor.Destination{}
	epIds := []string{}
	for i := range endpoints {
		signkey, err := encryption.DecryptAny(uc.conf.Infrastructure.Secrets.Cipher, endpoints[i].SecretKey)
		if err != nil {
			uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", err, "endpoint", utils.Stringify(endpoints[i]))
			continue
		}

		epIds = append(epIds, endpoints[i].Id)
		destinations[endpoints[i].Id] = &conductor.Destination{
			Endpoint: endpoints[i],
			SignKey:  signkey,
			Routes:   []*dbentities.Route{},
		}
	}

	if len(endpoints) > 0 && len(epIds) == 0 {
		return nil, ErrSchedulerArrange
	}

	routes, err := uc.getRoutesOfEndpoint(ctx, epIds)
	if err != nil {
		return nil, err
	}
	for i := range endpoints {
		destinations[endpoints[i].Id].Routes = routes[endpoints[i].Id]
	}

	return destinations, nil
}

func (uc *scheduler) getEndpointOfApp(ctx context.Context, appId string) ([]*dbentities.Endpoint, error) {
	var endpoints []*dbentities.Endpoint
	err := uc.orm.WithContext(ctx).
		Model(&dbentities.Endpoint{}).
		Where("app_id = ?", appId).
		Find(&endpoints).Error
	if err != nil {
		return nil, err
	}

	if len(endpoints) == 0 {
		uc.logger.Warnw("SDK.SCHEDULER.ARRANGE.NO_ENDPOINT.ERROR", "app_id", appId)
	}

	return endpoints, nil
}

func (uc *scheduler) getRoutesOfEndpoint(ctx context.Context, epIds []string) (map[string][]*dbentities.Route, error) {
	var routes []*dbentities.Route
	err := uc.orm.WithContext(ctx).
		Model(&dbentities.Route{}).
		Where("ep_id IN (?)", epIds).
		Order("ep_id DESC, exclusionary DESC, priority DESC").
		Find(&routes).Error
	if err != nil {
		return nil, err
	}

	maps := map[string][]*dbentities.Route{}
	for i := range routes {
		maps[routes[i].EpId] = append(maps[routes[i].EpId], routes[i])
	}

	for _, id := range epIds {
		if _, exist := maps[id]; !exist {
			uc.logger.Warnw("SDK.SCHEDULER.ARRANGE.NO_ROUTE.ERROR", "ep_id", id)
		}
	}

	return maps, nil
}

type SchedulerArrangeIn struct {
	Messages map[string]*dsentities.Message
}

func (in *SchedulerArrangeIn) Validate() error {
	return validator.Validate(validator.MapRequired("messages", in.Messages))
}

type SchedulerArrangeOut struct {
	Success []string
	Error   map[string]error
}
