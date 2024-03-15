package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/cipher/encryption"
	"github.com/kanthorlabs/common/safe"
	stmentities "github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/kanthor/internal/conductor"
	"github.com/kanthorlabs/kanthor/internal/constants"
	dbentities "github.com/kanthorlabs/kanthor/internal/database/entities"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
	"github.com/kanthorlabs/kanthor/internal/transformation"
	"github.com/sourcegraph/conc"
)

var ErrSchedulerArrange = errors.New("SDK.SCHDULER.ARRANGE.ERROR")

func (uc *scheduler) Arrange(ctx context.Context, in *SchedulerArrangeIn) (*SchedulerArrangeOut, error) {
	ok := &safe.Slice[string]{}
	ko := &safe.Map[error]{}

	refs := &safe.Map[string]{}

	// arrange requests to events
	events := &safe.Map[*stmentities.Event]{}
	var wg conc.WaitGroup
	for refId := range in.Messages {
		msg := in.Messages[refId]
		wg.Go(func() {
			reqs, err := uc.arrange(ctx, msg)
			if err != nil {
				uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", err, "message", utils.Stringify(msg))
				ko.Set(refId, err)
				return
			}

			for _, req := range reqs {
				// if any message request got error, reject the whole message
				// we don't accept partial success
				if _, exist := ko.Get(refId); exist {
					continue
				}

				event, err := transformation.EventFromRequest(req, constants.SubjectRequestSchedule)
				if err != nil {
					ko.Set(refId, err)
					uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", err, "request", utils.Stringify(req))
					continue
				}
				// one message can produce multiple events, so if you set event key by refId, it will be overwritten
				// then we create a direction of message -> refId -> epId -> event
				refs.Set(req.EpId, refId)
				events.Set(req.EpId, event)
			}
		})
	}
	wg.Wait()

	publisher, err := uc.infra.Streaming().Publisher(constants.PublisherRequest)
	if err != nil {
		uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", err.Error())
		return nil, ErrSchedulerArrange
	}

	errs := publisher.Pub(ctx, events.Data())
	if len(errs) > 0 {
		for epId, err := range errs {
			event, _ := events.Get(epId)
			uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", utils.Stringify(err), "event", utils.Stringify(event))

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
	epIds := make([]string, len(endpoints))
	for i := range endpoints {
		epIds[i] = endpoints[i].Id

		signkey, err := encryption.DecryptAny(uc.conf.Infrastructure.Secrets.Cipher, endpoints[i].SecretKey)
		if err != nil {
			uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", err, "endpoint", utils.Stringify(endpoints[i]))
			continue
		}

		destinations[endpoints[i].Id] = &conductor.Destination{
			Endpoint: endpoints[i],
			SignKey:  signkey,
			Routes:   []*dbentities.Route{},
		}
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
	return maps, nil
}

type SchedulerArrangeIn struct {
	Messages map[string]*dsentities.Message
}

type SchedulerArrangeOut struct {
	Success []string
	Error   map[string]error
}
