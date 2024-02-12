package routing

import (
	"fmt"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/timer"
	"github.com/kanthorlabs/kanthor/internal/entities"
)

func NewRequest(
	timer timer.Timer,
	msg *entities.Message,
	ep *entities.Endpoint,
	epr *entities.EndpointRule,
) *entities.Request {
	// construct request
	req := &entities.Request{
		MsgId:    msg.Id,
		Tier:     msg.Tier,
		AppId:    msg.AppId,
		Type:     msg.Type,
		EpId:     ep.Id,
		Metadata: entities.Metadata{},
		Headers:  entities.Header{},
		Body:     msg.Body,
		Uri:      ep.Uri,
		Method:   ep.Method,
	}

	// must use merge function otherwise you will edit the original data
	req.Headers.Merge(msg.Headers)
	req.Metadata.Merge(msg.Metadata)
	req.Id = idx.New(entities.IdNsReq)
	req.SetTS(timer.Now())

	req.Metadata.Set(entities.MetaEprId, epr.Id)

	req.Headers.Set("User-Agent", fmt.Sprintf("Kanthor/%s", project.GetVersion()))
	req.Headers.Set(entities.HeaderIdempotencyKey, msg.Id)

	return req
}
