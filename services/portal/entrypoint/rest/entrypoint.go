package rest

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/kanthor/database"
	"github.com/kanthorlabs/kanthor/datastore"
	"github.com/kanthorlabs/kanthor/gateway/gin/middlewares"
	"github.com/kanthorlabs/kanthor/infrastructure"
	"github.com/kanthorlabs/kanthor/openapi"
	"github.com/kanthorlabs/kanthor/patterns"
	"github.com/kanthorlabs/kanthor/services/portal/config"
	"github.com/kanthorlabs/kanthor/services/portal/usecase"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	infra *infrastructure.Infrastructure,
	db database.Database,
	ds datastore.Datastore,
	uc usecase.Portal,
) patterns.Runnable {
	logger = logger.With("service", "portal", "entrypoint", "rest")
	return &portal{
		conf:   conf,
		logger: logger,
		infra:  infra,
		db:     db,
		ds:     ds,
		uc:     uc,
	}
}

type portal struct {
	conf   *config.Config
	logger logging.Logger
	infra  *infrastructure.Infrastructure
	db     database.Database
	ds     datastore.Datastore
	uc     usecase.Portal

	server *http.Server

	mu     sync.Mutex
	status int
}

func (service *portal) Start(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status == patterns.StatusStarted {
		return ErrAlreadyStarted
	}

	if err := service.conf.Validate(); err != nil {
		return err
	}

	if err := service.db.Connect(ctx); err != nil {
		return err
	}

	if err := service.ds.Connect(ctx); err != nil {
		return err
	}

	if err := service.infra.Connect(ctx); err != nil {
		return err
	}

	router, err := service.router()
	if err != nil {
		return err
	}
	service.server = &http.Server{
		Addr:    service.conf.Gateway.Addr,
		Handler: router,
	}

	service.status = patterns.StatusStarted
	service.logger.Info("started")
	return nil
}

func (service *portal) router() (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.UseCors(service.conf.Gateway.Origins))
	// system routes
	RegisterHealthcheck(router, service)

	swagger := router.Group("/swagger")
	{
		swagger.GET("/*any", ginswagger.WrapHandler(
			swaggerfiles.Handler,
			ginswagger.PersistAuthorization(true),
			ginswagger.InstanceName(openapi.SwaggerInfoPortal.InfoInstanceName),
		))
	}

	api := router.Group("/api")
	{
		api.Use(middlewares.UseStartup(&service.conf.Gateway))
		api.Use(middlewares.UseTracing(project.Name("portal_rest")))
		api.Use(middlewares.UseIdempotency(service.logger, service.infra.Idempotency, project.IsDev()))
		api.Use(middlewares.UseAuthn(service.infra.Authenticator, service.infra.Authenticator.Engines()[0]))

		RegisterAnalyticsRoutes(api.Group("/analytics"), service)
		RegisterAccountRoutes(api.Group("/account"), service)
		RegisterWorkspaceRoutes(api.Group("/workspace"), service)
		RegisterWorkspaceCredentialsRoutes(api.Group("/credentials"), service)
		RegisterApplicationRoutes(api.Group("/application"), service)
		RegisterEndpointRoutes(api.Group("/endpoint"), service)
	}

	return router, nil
}

func (service *portal) Stop(ctx context.Context) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.status != patterns.StatusStarted {
		return ErrNotStarted
	}
	service.status = patterns.StatusStopped
	service.logger.Info("stopped")

	var returning error
	if err := service.server.Shutdown(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.infra.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.db.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	if err := service.ds.Disconnect(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	return returning
}

func (service *portal) Run(ctx context.Context) error {
	service.logger.Infow("running", "addr", service.conf.Gateway.Addr)
	if err := service.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
