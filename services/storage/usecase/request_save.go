package usecase

import (
	"context"

	"github.com/kanthorlabs/common/validator"
	dsentities "github.com/kanthorlabs/kanthor/internal/datastore/entities"
)

func (uc *request) Save(ctx context.Context, in *SaveRequestIn) (*SaveRequestOut, error) {
	var docs []*dsentities.Request
	for refId := range in.Requests {
		docs = append(docs, in.Requests[refId])
	}

	err := uc.repos.Request().Save(ctx, docs)
	if err != nil {
		out := &SaveRequestOut{Error: make(map[string]error, len(in.Requests))}
		for refId := range in.Requests {
			out.Error[refId] = err
		}
		return out, err
	}

	out := &SaveRequestOut{Success: make([]string, len(in.Requests))}
	for refId := range in.Requests {
		out.Success = append(out.Success, refId)
	}
	return out, err
}

type SaveRequestIn struct {
	Requests map[string]*dsentities.Request
}

func (in *SaveRequestIn) Validate() error {
	return validator.Validate(validator.MapRequired("requests", in.Requests))
}

type SaveRequestOut struct {
	Success []string
	Error   map[string]error
}
