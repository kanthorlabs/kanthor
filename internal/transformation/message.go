package transformation

import (
	"encoding/json"

	"github.com/kanthorlabs/common/project"
	stmentities "github.com/kanthorlabs/common/streaming/entities"

	"github.com/kanthorlabs/kanthor/internal/constants"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

func EventFromMessage(msg *dsentities.Message, subject string) (*stmentities.Event, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	event := &stmentities.Event{
		Subject: project.Subject(subject),
		Id:      msg.Id,
		Data:    data,
		Metadata: map[string]string{
			constants.MetadataProjectVersion: project.GetVersion(),
		},
	}

	return event, nil
}

func EventToMessage(event *stmentities.Event) (*dsentities.Message, error) {
	msg := &dsentities.Message{}
	if err := json.Unmarshal(event.Data, msg); err != nil {
		return nil, err
	}

	if err := msg.Validate(); err != nil {
		return nil, err
	}

	return msg, nil
}
