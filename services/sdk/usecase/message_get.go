package usecase

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	dbentities "github.com/kanthorlabs/kanthor/internal/database/entities"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

var ErrMessageGet = errors.New("SDK.MESSAGE.GET.ERROR")

func (uc *message) Get(ctx context.Context, in *MessageGetIn) (*MessageGetOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	out := &MessageGetOut{Endpoints: make(map[string]*MessageEndpoint)}

	message, err := uc.findMessage(ctx, in)
	if err != nil {
		return nil, err
	}
	// no message was found; return immediately because there is absolutely no existing request and response
	if message == nil {
		return out, nil
	}
	out.Message = message

	endpoints, err := uc.findEndpoints(ctx, message)
	if err != nil {
		return nil, err
	}
	out.Endpoints = endpoints

	return out, nil
}

func (uc *message) findMessage(ctx context.Context, in *MessageGetIn) (*dsentities.Message, error) {
	var pks []dsentities.MessagePk

	var apps []*dbentities.Application
	// @TODO: use channel with a chunk size for querying instead of querying all applications at once
	base := uc.orm.WithContext(ctx).Model(&dbentities.Application{}).Where("ws_id = ?", in.WsId)
	if err := base.Find(&apps).Error; err != nil {
		uc.logger.Errorw(ErrMessageGet.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrMessageGet
	}
	if len(apps) == 0 {
		return &dsentities.Message{}, nil
	}

	for i := range apps {
		pks = append(pks, dsentities.MessagePk{AppId: apps[i].Id, Id: in.Id})
	}

	messages, err := uc.repos.Message().Get(ctx, pks)
	if err != nil {
		uc.logger.Errorw(ErrMessageGet.Error(), "error", err.Error(), "message_pks", utils.Stringify(pks))
		return nil, ErrMessageGet
	}

	if len(messages) == 0 {
		return &dsentities.Message{}, nil
	}

	// we are expecting only receive one message of one app
	// because message is only post to ONE app
	return messages[0], nil
}

func (uc *message) findEndpoints(ctx context.Context, message *dsentities.Message) (map[string]*MessageEndpoint, error) {
	msgeps := make(map[string]*MessageEndpoint)

	// @TODO: use channel with a chunk size for querying instead of querying all endpoints at once
	var endpoints []*dbentities.Endpoint
	err := uc.orm.WithContext(ctx).Model(&dbentities.Endpoint{}).Where("app_id = ?", message.AppId).Find(&endpoints).Error
	if err != nil {
		uc.logger.Errorw(ErrMessageGet.Error(), "error", err.Error(), "message", utils.Stringify(message))
		return nil, ErrMessageGet
	}
	for i := range endpoints {
		msgeps[endpoints[i].Id] = &MessageEndpoint{Endpoint: endpoints[i]}
	}

	requests, err := uc.findRequests(ctx, message, endpoints)
	if err != nil {
		return nil, err
	}
	for _, request := range requests {
		msgeps[request.EpId].Requests = append(msgeps[request.EpId].Requests, request)
	}

	responses, err := uc.findResponses(ctx, message, endpoints)
	if err != nil {
		return nil, err
	}
	for _, response := range responses {
		msgeps[response.EpId].Responses = append(msgeps[response.EpId].Responses, response)
	}

	return msgeps, nil
}

func (uc *message) findRequests(ctx context.Context, message *dsentities.Message, eps []*dbentities.Endpoint) ([]*dsentities.Request, error) {
	var pks []dsentities.RequestPk
	for i := range eps {
		pks = append(pks, dsentities.RequestPk{EpId: eps[i].Id, MsgId: message.Id})
	}

	requests, err := uc.repos.Request().Get(ctx, pks)
	if err != nil {
		uc.logger.Errorw(ErrMessageGet.Error(), "error", err.Error(), "request_pks", utils.Stringify(pks))
		return nil, ErrMessageGet
	}

	return requests, nil
}

func (uc *message) findResponses(ctx context.Context, message *dsentities.Message, eps []*dbentities.Endpoint) ([]*dsentities.Response, error) {
	var pks []dsentities.ResponsePk
	for i := range eps {
		pks = append(pks, dsentities.ResponsePk{EpId: eps[i].Id, MsgId: message.Id})
	}

	responses, err := uc.repos.Response().Get(ctx, pks)
	if err != nil {
		uc.logger.Errorw(ErrMessageGet.Error(), "error", err.Error(), "response_pks", utils.Stringify(pks))
		return nil, ErrMessageGet
	}

	return responses, nil
}

type MessageGetIn struct {
	WsId string
	Id   string
}

func (in *MessageGetIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.MESSAGE.GET.IN.WS_ID", in.WsId, dbentities.IdNsWs),
		validator.StringStartsWith("SDK.MESSAGE.GET.IN.ID", in.Id, dsentities.IdNsMsg),
	)
}

type MessageGetOut struct {
	Message   *dsentities.Message
	Endpoints map[string]*MessageEndpoint
}

type MessageEndpoint struct {
	Endpoint  *dbentities.Endpoint
	Requests  []*dsentities.Request
	Responses []*dsentities.Response
}
