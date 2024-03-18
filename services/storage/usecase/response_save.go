package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

func (uc *response) Save(ctx context.Context, in *SaveResponseIn) (*SaveResponseOut, error) {
	var docs []*dsentities.Response
	for refId := range in.Responses {
		docs = append(docs, in.Responses[refId])
	}

	err := uc.repos.Response().Save(ctx, docs)
	if err != nil {
		out := &SaveResponseOut{Error: make(map[string]error, len(in.Responses))}
		for refId := range in.Responses {
			out.Error[refId] = err
		}
		return out, err
	}

	out := &SaveResponseOut{Success: make([]string, len(in.Responses))}
	for refId := range in.Responses {
		out.Success = append(out.Success, refId)
	}
	return out, err
}

type SaveResponseIn struct {
	Responses map[string]*dsentities.Response
}

func (in *SaveResponseIn) Validate() error {
	return validator.Validate(validator.MapRequired("requests", in.Responses))
}

type SaveResponseOut struct {
	Success []string
	Error   map[string]error
}
