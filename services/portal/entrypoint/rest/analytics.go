package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/kanthorlabs/kanthor/gateway/httpx/middlewares"
)

func RegisterAnalyticsRoutes(router gin.IRoutes, service *portal) {
	router = router.Use(middlewares.UseWorkspace(RegisterWorkspaceResolver(service.uc)))

	router.GET("overview", UseAnalyticsGetOverview(service))
}
