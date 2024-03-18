package transformation

import (
	"encoding/json"

	"github.com/kanthorlabs/common/project"
	stmentities "github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/kanthor/internal/constants"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

func EventFromRequest(req *dsentities.Request, subject string) (*stmentities.Event, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	event := &stmentities.Event{
		Subject: project.Subject(subject),
		Id:      req.Id,
		Data:    data,
		Metadata: map[string]string{
			constants.MetadataProjectVersion: project.GetVersion(),
		},
	}

	return event, nil
}

func EventToRequest(event *stmentities.Event) (*dsentities.Request, error) {
	req := &dsentities.Request{}
	if err := json.Unmarshal(event.Data, req); err != nil {
		return nil, err
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return req, nil
}
