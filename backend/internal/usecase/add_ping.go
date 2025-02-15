package usecase

import (
	"context"

	"backend/internal/domain"
	"backend/internal/infrastructure"
)

type AddPingUseCase struct {
	PingRepo infrastructure.PingRepository
}

func NewPingUseCase(pingRepo infrastructure.PingRepository) *AddPingUseCase {
	return &AddPingUseCase{
		PingRepo: pingRepo,
	}
}

func (uc *AddPingUseCase) CreatePing(ctx context.Context, p *domain.Ping) error {
	return uc.PingRepo.Create(ctx, p)
}
