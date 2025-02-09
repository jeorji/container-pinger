package domain

import "context"

type ContainerRepository interface {
	UpsertByID(ctx context.Context, c *Container) error
	GetAll(ctx context.Context) ([]Container, error)
}

type PingRepository interface {
	Create(ctx context.Context, p *Ping) error
	GetLastPingForContainer(ctx context.Context, p *Container) ([]Ping, error)
}
