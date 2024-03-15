package conductor

import (
	"fmt"
	"time"

	"github.com/kanthorlabs/common/cipher/signature"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/safe"
	dbentities "github.com/kanthorlabs/kanthor/internal/database/entities"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

var MetaRouteId = "kanthor.rt.id"

func Request(
	msg *dsentities.Message,
	ep *dbentities.Endpoint,
	rt *dbentities.Route,
	now time.Time,
	signkey string,
) *dsentities.Request {
	req := &dsentities.Request{
		EpId: ep.Id,

		MsgId:    msg.Id,
		Tier:     msg.Tier,
		AppId:    msg.AppId,
		Type:     msg.Type,
		Metadata: &safe.Metadata{},
		Body:     msg.Body,

		Method:  ep.Method,
		Uri:     ep.Uri,
		Headers: &safe.Metadata{},
	}

	req.SetId()
	req.SetTimeseries(now)
	req.Metadata.Merge(msg.Metadata)
	req.Metadata.Set(MetaRouteId, rt.Id)
	req.Headers.Set("User-Agent", fmt.Sprintf("Kanthor/%s", project.GetVersion()))
	req.Headers.Set("Idempotency-Key", msg.Id)

	req.Headers.Set("Webhook-Id", req.Id)
	req.Headers.Set("Webhook-Timestamp", fmt.Sprintf("%d", req.CreatedAt))
	req.Headers.Set("Webhook-Signature", signature.Sign(signkey, fmt.Sprintf("%s.%d.%s", req.Id, req.CreatedAt, req.Body)))
	req.Headers.Set("Webhook-Type", req.Type)

	return req
}
