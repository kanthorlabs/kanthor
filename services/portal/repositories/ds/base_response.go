package ds

import (
	"context"

	"github.com/kanthorlabs/kanthor/internal/entities"
)

type MessageResponsetMaps struct {
	Maps    map[string][]entities.Response
	Success map[string]string
}

type Response interface {
	ListMessages(ctx context.Context, epId string, msgIds []string) (*MessageResponsetMaps, error)
	GetMessages(ctx context.Context, epId string, msgIds []string) (*MessageResponsetMaps, error)
}
