package restapi

import (
	"backend/internal/domain"
	"backend/internal/usecase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ContainerStatsUC  *usecase.ContainerStatsUseCase
	ContainerUpdateUC *usecase.ContainerUpdateUseCase
	PingUC            *usecase.AddPingUseCase
}

func NewHandler(statsUC *usecase.ContainerStatsUseCase, updateUC *usecase.ContainerUpdateUseCase, pingUC *usecase.AddPingUseCase) *Handler {
	return &Handler{
		ContainerStatsUC:  statsUC,
		ContainerUpdateUC: updateUC,
		PingUC:            pingUC,
	}
}

func (h *Handler) GetContainers(c *gin.Context) {
	containers, err := h.ContainerStatsUC.GetAllStats(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, containers)
}

func (h *Handler) CreateOrUpdateContainer(c *gin.Context) {
	containerID := c.Param("id")

	var dto domain.Container
	dto.ID = containerID

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ContainerUpdateUC.CreateOrUpdateContainerByID(c, &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) AddPingRecord(c *gin.Context) {
	containerID := c.Param("id")

	var dto domain.Ping
	dto.ContainerID = containerID

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("ERR: %s\n", err)
		return
	}

	err := h.PingUC.CreatePing(c, &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("ERR: %s\n", err)
		return
	}

	c.Status(http.StatusCreated)
}
