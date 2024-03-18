package transformation

import (
	"encoding/json"

	"github.com/kanthorlabs/common/project"
	stmentities "github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/kanthor/internal/constants"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

func EventFromResponse(res *dsentities.Response, subject string) (*stmentities.Event, error) {
	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	event := &stmentities.Event{
		Subject: project.Subject(subject),
		Id:      res.Id,
		Data:    data,
		Metadata: map[string]string{
			constants.MetadataProjectVersion: project.GetVersion(),
		},
	}

	return event, nil
}

func EventToResponse(event *stmentities.Event) (*dsentities.Response, error) {
	res := &dsentities.Response{}
	if err := json.Unmarshal(event.Data, res); err != nil {
		return nil, err
	}

	if err := res.Validate(); err != nil {
		return nil, err
	}

	return res, nil
}
