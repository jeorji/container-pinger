package domain

type ContainerRepository interface {
	ListContainers(all bool) ([]ContainerInfo, error)
	InspectContainer(containerID string) (ContainerInfo, []string, error)
}

type BackendRepository interface {
	UpdateContainer(containerInfo ContainerInfo) error
	PostPing(pingRes PingResult) error
}

type PingService interface {
	Ping(ip string) (status string, latency float32, err error)
}
