package infrastructure

import (
	"backend/internal/domain"
	"context"
)

type ContainerRepository interface {
	UpsertByID(ctx context.Context, c *domain.Container) error
	GetAll(ctx context.Context) ([]domain.Container, error)
}

type PingRepository interface {
	Create(ctx context.Context, p *domain.Ping) error
	GetLastPingForContainer(ctx context.Context, p *domain.Container) ([]domain.Ping, error)
}
