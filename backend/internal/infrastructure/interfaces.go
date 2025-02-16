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

// Mocks

type MockContainerRepo struct {
	GetAllFunc func(ctx context.Context) ([]domain.Container, error)
	UpsertByIDFunc func(ctx context.Context, c *domain.Container) error
}

func (m *MockContainerRepo) GetAll(ctx context.Context) ([]domain.Container, error) {
	return m.GetAllFunc(ctx)
}
func (m *MockContainerRepo) UpsertByID(ctx context.Context, c *domain.Container) error {
	return m.UpsertByIDFunc(ctx, c)
}

type MockPingRepo struct {
	GetLastPingForContainerFunc func(ctx context.Context, c *domain.Container) ([]domain.Ping, error)
	CreateFunc func(ctx context.Context, p *domain.Ping) error
}

func (m *MockPingRepo) GetLastPingForContainer(ctx context.Context, c *domain.Container) ([]domain.Ping, error) {
	return m.GetLastPingForContainerFunc(ctx, c)
}
func (m *MockPingRepo) Create(ctx context.Context, p *domain.Ping) error {
	return m.CreateFunc(ctx, p)
}
