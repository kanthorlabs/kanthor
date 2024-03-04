package api

import (
	"github.com/kanthorlabs/common/gateway/httpx"
	"github.com/sourcegraph/conc/pool"
)

func (service *portal) httpx() error {
	handler, err := httpx.New(&service.conf.Portal.Gateway, service.logger)
	if err != nil {
		return err
	}

	handler.Get("/healthz/readiness", httpx.UseHealthz(func() error {
		p := pool.New().WithErrors()
		p.Go(func() error {
			return service.infra.Readiness()
		})
		return p.Wait()
	}))

	handler.Get("/healthz/liveness", httpx.UseHealthz(func() error {
		p := pool.New().WithErrors()
		p.Go(func() error {
			return service.infra.Liveness()
		})
		return p.Wait()
	}))

	return service.server.UseHttpx(handler)
}
