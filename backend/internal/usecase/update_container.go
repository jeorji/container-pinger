
package usecase

import (
	"context"

	"backend/internal/domain"
)

type ContainerUpdateUseCase struct {
	ContainerRepo  domain.ContainerRepository
}

func NewContainerUpdateUseCase(
	containerRepo domain.ContainerRepository,
) *ContainerUpdateUseCase {
	return &ContainerUpdateUseCase{
        ContainerRepo: containerRepo,
	}
}

func (uc *ContainerUpdateUseCase) CreateOrUpdateContainerByID(ctx context.Context, c *domain.Container) error {
	err := uc.ContainerRepo.UpsertByID(ctx, c)
	return err
}
