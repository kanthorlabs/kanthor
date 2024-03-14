package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kanthorlabs/common/cache"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/safe"
	stmentities "github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"github.com/kanthorlabs/kanthor/internal/constants"
	dbentities "github.com/kanthorlabs/kanthor/internal/database/entities"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

var AppTierCacheTTL = time.Hour * 24
var ErrMessageCreate = errors.New("SDK.MESSAGE.CREATE.ERROR")

func (uc *message) Create(ctx context.Context, in *MessageCreateIn) (*MessageCreateOut, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	tier, err := uc.getAppTier(ctx, in.AppId)
	if err != nil {
		uc.logger.Errorw(ErrMessageCreate.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrMessageCreate
	}

	msg := &dsentities.Message{
		Tier:     tier,
		AppId:    in.AppId,
		Tag:      in.Tag,
		Metadata: &safe.Metadata{},
		Body:     in.Body,
	}
	msg.SetId()
	msg.SetTimeseries(uc.watch.Now())
	msg.Metadata.Set(constants.MetadataProjectVersion, project.GetVersion())

	// no way to get error here
	data, err := json.Marshal(msg)
	if err != nil {
		uc.logger.Errorw(ErrMessageCreate.Error(), "error", err.Error(), "in", utils.Stringify(in))
		return nil, ErrMessageCreate
	}

	publisher, err := uc.infra.Streaming().Publisher(constants.PublisherMessage)
	if err != nil {
		uc.logger.Errorw(ErrMessageCreate.Error(), "error", err.Error(), "msg", string(data))
		return nil, ErrMessageCreate
	}

	event := &stmentities.Event{
		Subject: project.Subject(constants.SubjectMessageCreate),
		Id:      msg.Id,
		Data:    data,
		Metadata: map[string]string{
			constants.MetadataProjectVersion: project.GetVersion(),
		},
	}
	errs := publisher.Pub(ctx, map[string]*stmentities.Event{event.Id: event})
	if len(errs) > 0 {
		uc.logger.Errorw(ErrMessageCreate.Error(), "error", errs[event.Id], "event", utils.Stringify(event))
		return nil, ErrMessageCreate
	}

	return &MessageCreateOut{Id: msg.Id, CreatedAt: msg.CreatedAt}, nil
}

func (uc *message) getAppTier(ctx context.Context, appId string) (string, error) {
	tier, err := cache.GetOrSet(
		uc.infra.Cache(), ctx,
		constants.CacheKeyAppTier(appId),
		AppTierCacheTTL,
		func() (*string, error) {
			sql := fmt.Sprintf(
				"SELECT tier FROM %s JOIN %s ON %s.ws_id = %s.id WHERE %s.id = ?",
				dbentities.TableWs, dbentities.TableApp, dbentities.TableApp, dbentities.TableWs, dbentities.TableApp,
			)

			var tier string
			return &tier, uc.orm.Raw(sql, appId).Scan(&tier).Error
		})

	if err != nil {
		return "", err
	}
	return *tier, nil
}

type MessageCreateIn struct {
	AppId string
	Tag   string
	Body  string
}

func (in *MessageCreateIn) Validate() error {
	return validator.Validate(
		validator.StringStartsWith("SDK.MESSAGE.CREATE.IN.APP_ID", in.AppId, dbentities.IdNsApp),
		validator.StringAlphaNumericUnderscoreDot("SDK.MESSAGE.CREATE.IN.TAG", in.Tag),
		validator.StringRequired("SDK.MESSAGE.CREATE.IN.BODY", in.Body),
	)
}

type MessageCreateOut struct {
	Id        string
	CreatedAt int64
}
