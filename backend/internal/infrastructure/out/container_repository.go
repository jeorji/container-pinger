package repository

import (
	"backend/internal/domain"
	"backend/internal/infrastructure"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type containerPG struct {
	db *pgxpool.Pool
}
func NewContainerPGRepo(db *pgxpool.Pool) infrastructure.ContainerRepository {
	return &containerPG{db: db}
}

func (c *containerPG) GetAll(ctx context.Context) ([]domain.Container, error) {
    query := `
        SELECT id, name, image, state, status, created_at
        FROM containers
    `

    rows, err := c.db.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var containers []domain.Container
    for rows.Next() {
        var container domain.Container
        err := rows.Scan(&container.ID, &container.Name, &container.Image, &container.State, &container.Status, &container.CreatedAt)
        if err != nil {
            return nil, err
        }
        containers = append(containers, container)
    }

    return containers, nil
}

func (c *containerPG) UpsertByID(ctx context.Context, cont *domain.Container) error {
	query := `
        INSERT INTO containers (id, name, image, state, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (id) DO UPDATE SET
            name = EXCLUDED.name,
            state = EXCLUDED.state,
            status = EXCLUDED.status;
    `

	_, err := c.db.Exec(ctx, query, cont.ID, cont.Name, cont.Image, cont.State, cont.Status, cont.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to upsert container: %w", err)
	}

	return nil
}
