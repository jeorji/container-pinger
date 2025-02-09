package restapi

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *Handler) {
	apiGroup := r.Group("/api")

	containersGroup := apiGroup.Group("/containers")
	{
		containersGroup.GET("", h.GetContainers)              // GET /api/containers
		containersGroup.PUT(":id", h.CreateOrUpdateContainer) // PUT /api/containers/:id
		containersGroup.POST(":id/pings", h.AddPingRecord)    // POST /api/containers/:id/pings
	}
}
