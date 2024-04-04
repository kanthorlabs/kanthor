package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanthorlabs/common/cache"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/safe"
	stmentities "github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/constants"
	dbentities "github.com/kanthorlabs/kanthor/internal/database/entities"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
	"github.com/kanthorlabs/kanthor/internal/transformation"
	"github.com/kanthorlabs/kanthor/services/sdk/caching"
)

var ErrMessageCreate = errors.New("SDK.MESSAGE.CREATE.ERROR")

func (uc *message) Create(ctx context.Context, in *MessageCreateIn) (*MessageCreateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	app, err := uc.getApp(ctx, in.AppId)
	if err != nil {
		uc.logger.Errorw(ErrMessageCreate.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrMessageCreate
	}

	msg := &dsentities.Message{
		Tier:     app.Tier,
		AppId:    app.Id,
		Type:     in.Type,
		Body:     in.Body,
		Headers:  &safe.Metadata{},
		Metadata: &safe.Metadata{},
	}
	msg.SetId()
	msg.SetTimeseries(uc.watch.Now())
	msg.Metadata.Set(constants.MetadataProjectVersion, project.GetVersion())

	// message is the most important entity in the system
	// without it we cannot process further process like recovery or attempt
	// so make sure we save it before do anything else
	if err := uc.repos.Message().Save(ctx, []*dsentities.Message{msg}); err != nil {
		uc.logger.Errorw(ErrMessageCreate.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrMessageCreate
	}

	event, err := transformation.EventFromMessage(msg, constants.MessageCreateSubject)
	if err != nil {
		uc.logger.Errorw(ErrMessageCreate.Error(), "error", err.Error(), "msg", utils.Stringify(msg))
		return nil, ErrMessageCreate
	}

	publisher, err := uc.infra.Streaming().Publisher(constants.MessagePublisher)
	if err != nil {
		uc.logger.Errorw(ErrMessageCreate.Error(), "error", err.Error(), "msg", utils.Stringify(msg))
		return nil, ErrMessageCreate
	}

	errs := publisher.Pub(ctx, map[string]*stmentities.Event{event.Id: event})
	if len(errs) > 0 {
		uc.logger.Errorw(ErrMessageCreate.Error(), "error", errs[event.Id], "event", event.String())
		return nil, ErrMessageCreate
	}

	return &MessageCreateOut{Id: msg.Id, CreatedAt: msg.CreatedAt}, nil
}

func (uc *message) getApp(ctx context.Context, appId string) (*MessageApp, error) {
	key, duration := caching.App(appId)
	return cache.GetOrSet(
		uc.infra.Cache(), ctx, key, duration,
		func() (*MessageApp, error) {
			sql := fmt.Sprintf(
				"SELECT %s.id, %s.tier FROM %s JOIN %s ON %s.id = %s.ws_id WHERE %s.id = ?",
				dbentities.TableApp, dbentities.TableWs,
				dbentities.TableApp,
				dbentities.TableWs, dbentities.TableWs, dbentities.TableApp,
				dbentities.TableApp,
			)

			app := MessageApp{}
			return &app, uc.orm.Raw(sql, appId).Scan(&app).Error
		})
}

type MessageCreateIn struct {
	WsId  string
	AppId string
	Type  string
	Body  string
}

func (in *MessageCreateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.MESSAGE.CREATE.IN.WS_ID", in.WsId, dbentities.IdNsWs),
		validator.StringStartsWith("SDK.MESSAGE.CREATE.IN.APP_ID", in.AppId, dbentities.IdNsApp),
		validator.StringAlphaNumericUnderscoreHyphenDot("SDK.MESSAGE.CREATE.IN.TYPE", in.Type),
		validator.StringRequired("SDK.MESSAGE.CREATE.IN.BODY", in.Body),
	)
}

type MessageCreateOut struct {
	Id        string
	CreatedAt int64
}

type MessageApp struct {
	Id   string
	Tier string
}
