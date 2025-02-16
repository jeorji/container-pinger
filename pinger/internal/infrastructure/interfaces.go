package infrastructure 

import "pinger/internal/domain"

type ContainerRepository interface {
	ListContainers(all bool) ([]domain.ContainerInfo, error)
	InspectContainer(containerID string) (domain.ContainerInfo, []string, error)
}

type BackendRepository interface {
	UpdateContainer(containerInfo domain.ContainerInfo) error
	PostPing(pingRes domain.PingResult) error
}

type PingService interface {
	Ping(ip string) (status string, latency float32, err error)
}
