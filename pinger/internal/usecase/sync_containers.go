package usecase

import (
	"context"
	"log"
	"sync"

	"pinger/internal/domain"
	"pinger/internal/infrastructure"
)

type ContainerSyncUseCase struct {
	DockerRepo  infrastructure.ContainerRepository
	BackendRepo infrastructure.BackendRepository
}

func NewContainerSyncUseCase(
	dockerRepo infrastructure.ContainerRepository,
	backendRepo infrastructure.BackendRepository,
) *ContainerSyncUseCase {
	return &ContainerSyncUseCase{
		DockerRepo:  dockerRepo,
		BackendRepo: backendRepo,
	}
}

func (uc *ContainerSyncUseCase) SyncContainers(ctx context.Context, all bool) error {
	containers, err := uc.DockerRepo.ListContainers(all)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, c := range containers {
		wg.Add(1)
		go func(c domain.ContainerInfo) {
			defer wg.Done()
			fullInfo, _, err := uc.DockerRepo.InspectContainer(c.ID)
			if err != nil {
				log.Printf("%s: %v", c.ID, err)
				return
			}
			if err := uc.BackendRepo.UpdateContainer(fullInfo); err != nil {
				log.Printf("%s: %v", c.ID, err)
			}
		}(c)
	}

	wg.Wait()
	return nil
}
