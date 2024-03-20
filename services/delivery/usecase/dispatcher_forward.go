package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/circuitbreaker"
	"github.com/kanthorlabs/common/safe"
	sndentities "github.com/kanthorlabs/common/sender/entities"
	stmentities "github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/constants"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
	"github.com/kanthorlabs/kanthor/internal/transformation"
	"github.com/sourcegraph/conc/pool"
)

var ErrDispatcherForward = errors.New("DELIVERY.DISPATCHER.FORWARD.ERROR")

func (uc *dispatcher) Fowrard(ctx context.Context, in *DispatcherForwardIn) (*DispatcherForwardOut, error) {
	ok := &safe.Slice[string]{}
	ko := &safe.Map[error]{}

	events := &safe.Map[*stmentities.Event]{}

	p := pool.New().WithContext(ctx)
	for id := range in.Requests {
		refId := id
		request := in.Requests[refId]

		p.Go(func(subctx context.Context) error {
			response, err := uc.forward(subctx, request)
			// the error indicates that the request get unrecoverable error
			if err != nil {
				uc.logger.Errorw(ErrDispatcherForward.Error(), "error", err.Error(), "request", utils.Stringify(request), "ref_id", refId)
				return nil
			}

			event, err := transformation.EventFromResponse(response, constants.ResponseNotifySubject)
			if err != nil {
				uc.logger.Errorw(ErrDispatcherForward.Error(), "error", err.Error(), "request", utils.Stringify(request), "response", utils.Stringify(response), "ref_id", refId)
				return nil
			}

			// can reuse refId here because one request is only produce one response
			events.Set(refId, event)
			// must not return an error here, we will handle it with the KO map
			return nil
		})
	}

	// got the context timeout error, should return the error to get it re-delivered
	if err := p.Wait(); err != nil {
		for id := range in.Requests {
			ko.Set(id, err)
		}
	}

	publisher, err := uc.infra.Streaming().Publisher(constants.ResponsePublisher)
	if err != nil {
		uc.logger.Errorw(ErrDispatcherForward.Error(), "error", err.Error())
		return nil, ErrDispatcherForward
	}

	errs := publisher.Pub(ctx, events.Data())
	if len(errs) > 0 {
		// because we use the refId as the key, we can obtain the refId from the error map
		for refId := range errs {
			event, _ := events.Get(refId)
			uc.logger.Errorw(ErrSchedulerArrange.Error(), "error", utils.Stringify(err), "event", event.String())

			ko.Set(refId, err)
		}
	}

	return &DispatcherForwardOut{Success: ok.Data(), Error: ko.Data()}, nil
}

func (uc *dispatcher) forward(ctx context.Context, request *dsentities.Request) (*dsentities.Response, error) {
	return circuitbreaker.Do[dsentities.Response](
		uc.infra.CircuitBreaker(),
		request.EpId,
		func() (any, error) {
			req := &sndentities.Request{
				Method:  request.Method,
				Headers: request.Headers.ToHttpHeader(),
				Uri:     request.Uri,
				Body:    []byte(request.Body),
			}
			res, err := uc.send(ctx, req)
			// got un-recovable error, should return the error to ignore the response
			if err != nil {
				return nil, err
			}

			response := &dsentities.Response{
				MsgId:    request.MsgId,
				Tier:     request.Tier,
				AppId:    request.AppId,
				Type:     request.Type,
				ReqId:    request.Id,
				EpId:     request.EpId,
				Headers:  &safe.Metadata{},
				Metadata: &safe.Metadata{},
				Status:   res.Status,
				Uri:      res.Uri,
				Body:     string(res.Body),
			}
			response.Headers.FromHttpHeader(res.Headers)
			response.Metadata.Merge(request.Metadata)
			response.SetId()
			response.SetTimeseries(uc.watch.Now())
			return response, nil
		},
		func(err error) error { return err },
	)
}

type DispatcherForwardIn struct {
	Requests map[string]*dsentities.Request
}

func (in *DispatcherForwardIn) Validate() error {
	return validator.Validate(validator.MapRequired("requests", in.Requests))
}

type DispatcherForwardOut struct {
	Success []string
	Error   map[string]error
}
