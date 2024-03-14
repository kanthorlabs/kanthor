package usecase

import (
	"context"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/services/sdk/config"
	"gorm.io/gorm"
)

type Endpoint interface {
	Create(ctx context.Context, in *EndpointCreateIn) (*EndpointCreateOut, error)
	Get(ctx context.Context, in *EndpointGetIn) (*EndpointGetOut, error)
	List(ctx context.Context, in *EndpointListIn) (*EndpointListOut, error)
	Update(ctx context.Context, in *EndpointUpdateIn) (*EndpointUpdateOut, error)
	Delete(ctx context.Context, in *EndpointDeleteIn) (*EndpointDeleteOut, error)

	// GetSecretOwn is a method to get the endpoint by its id and the workspace id
	GetOwn(ctx context.Context, in *EndpointGetOwnIn) (*EndpointGetOwnOut, error)
	// GetSecret return the endpoint with the decrypted secret key, so user can use it to verify the request to the endpoint
	GetSecret(ctx context.Context, in *EndpointGetSecretIn) (*EndpointGetSecretOut, error)
}

type endpoint struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	infra  infrastructure.Infrastructure
	orm    *gorm.DB
}
