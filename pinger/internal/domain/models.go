package domain

import "time"

type ContainerInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	State     string    `json:"state"`      // "running", "exited" и т.д.
	Status    string    `json:"status"`     // "Up 10 seconds" и т.д.
	CreatedAt time.Time `json:"created_at"`
}

type PingResult struct {
	ContainerID string    `json:"container_id"`
	IP          string    `json:"ip"`
	Status      string    `json:"status"`
	Latency     float32   `json:"latency"`
	PingTime    time.Time `json:"ping_time"`
}

