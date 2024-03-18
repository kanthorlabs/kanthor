package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

func (uc *message) Save(ctx context.Context, in *SaveMessageIn) (*SaveMessageOut, error) {
	var docs []*dsentities.Message
	for refId := range in.Messages {
		docs = append(docs, in.Messages[refId])
	}

	err := uc.repos.Message().Save(ctx, docs)
	if err != nil {
		out := &SaveMessageOut{Error: make(map[string]error, len(in.Messages))}
		for refId := range in.Messages {
			out.Error[refId] = err
		}
		return out, err
	}

	out := &SaveMessageOut{Success: make([]string, len(in.Messages))}
	for refId := range in.Messages {
		out.Success = append(out.Success, refId)
	}
	return out, err
}

type SaveMessageIn struct {
	Messages map[string]*dsentities.Message
}

func (in *SaveMessageIn) Validate() error {
	return validator.Validate(validator.MapRequired("messages", in.Messages))
}

type SaveMessageOut struct {
	Success []string
	Error   map[string]error
}
