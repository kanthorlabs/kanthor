package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gateway/httpx"
	"github.com/sourcegraph/conc/pool"
)

func RegisterHealthzRoutes(router chi.Router, service *portal) {
	router.Get("/healthz/readiness", httpx.UseHealthz(func() error {
		p := pool.New().WithErrors()
		p.Go(func() error {
			return service.infra.Readiness()
		})
		return p.Wait()
	}))

	router.Get("/healthz/liveness", httpx.UseHealthz(func() error {
		p := pool.New().WithErrors()
		p.Go(func() error {
			return service.infra.Liveness()
		})
		return p.Wait()
	}))
}
