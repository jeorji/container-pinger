package usecase

import (
	"context"

	"backend/internal/domain"
	"backend/internal/infrastructure"
)

type ContainerUpdateUseCase struct {
	ContainerRepo  infrastructure.ContainerRepository
}

func NewContainerUpdateUseCase(
	containerRepo infrastructure.ContainerRepository,
) *ContainerUpdateUseCase {
	return &ContainerUpdateUseCase{
        ContainerRepo: containerRepo,
	}
}

func (uc *ContainerUpdateUseCase) CreateOrUpdateContainerByID(ctx context.Context, c *domain.Container) error {
	err := uc.ContainerRepo.UpsertByID(ctx, c)
	return err
}
