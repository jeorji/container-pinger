package usecase

import (
	"context"

	"backend/internal/domain"
)

type ContainerStatsUseCase struct {
	ContainerRepo domain.ContainerRepository
	PingRepo      domain.PingRepository
}

func NewContainerStatsUseCase(containerRepo domain.ContainerRepository, pingRepo domain.PingRepository) *ContainerStatsUseCase {
	return &ContainerStatsUseCase{
		ContainerRepo: containerRepo,
		PingRepo:      pingRepo,
	}
}

func (uc *ContainerStatsUseCase) GetAllStats(ctx context.Context) ([]domain.ContainerStatsDTO, error) {
	containers, err := uc.ContainerRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var stats []domain.ContainerStatsDTO

	for _, container := range containers {
		pings, err := uc.PingRepo.GetLastPingForContainer(ctx, &container)
		if err != nil {
			return nil, err
		}

		pingStats := make([]domain.PingStatsDTO, len(pings))
		for i, ping := range pings {
			pingStats[i] = domain.PingStatsDTO{
				IP:             ping.IP,
				LastLatency:    ping.Latency,
				LastSuccessful: ping.PingTime,
			}
		}

		stats = append(stats, domain.ContainerStatsDTO{
			ID:        container.ID,
			Name:      container.Name,
			Image:     container.Image,
			State:     container.State,
			Status:    container.Status,
			CreatedAt: container.CreatedAt,
			PingStats: pingStats,
		})
	}

	return stats, nil
}
