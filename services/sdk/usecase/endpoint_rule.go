package usecase

import (
	"context"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"github.com/kanthorlabs/kanthor/services/sdk/repositories"
)

type EndpointRule interface {
	Create(ctx context.Context, in *EndpointRuleCreateIn) (*EndpointRuleCreateOut, error)
	Update(ctx context.Context, in *EndpointRuleUpdateIn) (*EndpointRuleUpdateOut, error)
	Delete(ctx context.Context, in *EndpointRuleDeleteIn) (*EndpointRuleDeleteOut, error)

	List(ctx context.Context, in *EndpointRuleListIn) (*EndpointRuleListOut, error)
	Get(ctx context.Context, in *EndpointRuleGetIn) (*EndpointRuleGetOut, error)
}

type endpointRule struct {
	conf         *config.Config
	logger       logging.Logger
	infra        *infrastructure.Infrastructure
	repositories repositories.Repositories
}
