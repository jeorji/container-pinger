package repository

import (
	"backend/internal/domain"
	"backend/internal/infrastructure"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pingPG struct {
	db *pgxpool.Pool
}

func NewPingPGRepo(db *pgxpool.Pool) infrastructure.PingRepository {
	return &pingPG{db: db}
}

func (r *pingPG) GetLastPingForContainer(ctx context.Context, p *domain.Container) ([]domain.Ping, error) {
	query := `
        SELECT DISTINCT ON (ip) ip, status, latency, ping_time
        FROM pings
        WHERE container_id = $1 AND status = 'success'
        ORDER BY ip, ping_time DESC
    `

	rows, err := r.db.Query(ctx, query, p.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pings []domain.Ping
	for rows.Next() {
		var ping domain.Ping
		if err := rows.Scan(&ping.IP, &ping.Status, &ping.Latency, &ping.PingTime); err != nil {
			return nil, err
		}
		ping.ContainerID = p.ID
		pings = append(pings, ping)
	}

	return pings, nil
}

func (r *pingPG) Create(ctx context.Context, p *domain.Ping) error {
	query := `
        INSERT INTO pings (container_id, ip, status, latency, ping_time)
        VALUES ($1, $2, $3, $4, $5);
    `

	_, err := r.db.Exec(ctx, query, p.ContainerID, p.IP, p.Status, p.Latency, p.PingTime)
	if err != nil {
		return fmt.Errorf("failed to insert ping: %w", err)
	}

	return nil
}
