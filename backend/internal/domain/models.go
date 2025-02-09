package domain

import "time"

type Container struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	State     string `json:"state"`
	Status    string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
type Ping struct {
	ContainerID string    `json:"container_id"`
	IP          string    `json:"ip"`
	Status      string    `json:"status"`
	Latency     float32   `json:"latency"`
	PingTime    time.Time `json:"ping_time"`
}


// TODO: move to usecase/dto
type ContainerStatsDTO struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Image     string         `json:"image"`
	State     string         `json:"state"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	PingStats []PingStatsDTO `json:"pings"`
}
type PingStatsDTO struct {
	IP             string    `json:"ip"`
	LastLatency    float32   `json:"last_ping_latency"`
	LastSuccessful time.Time `json:"last_ping_time"`
}
