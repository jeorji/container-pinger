package usecase

import (
	"context"

	"backend/internal/domain"
)

type AddPingUseCase struct {
	PingRepo domain.PingRepository
}

func NewPingUseCase(pingRepo domain.PingRepository) *AddPingUseCase {
	return &AddPingUseCase{
		PingRepo: pingRepo,
	}
}

func (uc *AddPingUseCase) CreatePing(ctx context.Context, p *domain.Ping) error {
	return uc.PingRepo.Create(ctx, p)
}
