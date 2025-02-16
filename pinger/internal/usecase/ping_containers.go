package usecase

import (
	"context"
	"log"
	"sync"
	"time"

	"pinger/internal/domain"
	"pinger/internal/infrastructure"
)

type PingContainersUseCase struct {
	DockerRepo  infrastructure.ContainerRepository
	BackendRepo infrastructure.BackendRepository
	PingService infrastructure.PingService
}

func NewPingContainersUseCase(
	dockerRepo infrastructure.ContainerRepository,
	backendRepo infrastructure.BackendRepository,
	pingService infrastructure.PingService,
) *PingContainersUseCase {
	return &PingContainersUseCase{
		DockerRepo:  dockerRepo,
		BackendRepo: backendRepo,
		PingService: pingService,
	}
}

func (uc *PingContainersUseCase) PingContainers(ctx context.Context, all bool) error {
	containers, err := uc.DockerRepo.ListContainers(all)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, c := range containers {
		wg.Add(1)
		go func(c domain.ContainerInfo) {
			defer wg.Done()
			_, ips, err := uc.DockerRepo.InspectContainer(c.ID)
			if err != nil {
				log.Printf("%s: %v", c.ID, err)
				return
			}

			for _, ip := range ips {
				if ip == "" {
					continue
				}
				status, latency, pingErr := uc.PingService.Ping(ip)
				if pingErr != nil {
					log.Printf("%s: %s %v", c.ID, ip, pingErr)
				}

				pingRes := domain.PingResult{
					ContainerID: c.ID,
					IP:          ip,
					Status:      status,
					Latency:     latency,
					PingTime:    time.Now(),
				}

				if err := uc.BackendRepo.PostPing(pingRes); err != nil {
					log.Printf("%s: %v", c.ID, err)
				}

			}
		}(c)
	}
	wg.Wait()
	return nil
}
