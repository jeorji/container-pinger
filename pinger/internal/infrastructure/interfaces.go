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

// Mocks

type MockContainerRepo struct {
	ListContainersFunc   func(all bool) ([]domain.ContainerInfo, error)
	InspectContainerFunc func(containerID string) (domain.ContainerInfo, []string, error)
}

func (m *MockContainerRepo) ListContainers(all bool) ([]domain.ContainerInfo, error) {
	return m.ListContainersFunc(all)
}
func (m *MockContainerRepo) InspectContainer(containerID string) (domain.ContainerInfo, []string, error) {
	return m.InspectContainerFunc(containerID)
}

type MockBackendRepo struct {
	UpdateContainerFunc func(containerInfo domain.ContainerInfo) error
	PostPingFunc        func(pingRes domain.PingResult) error
}

func (m *MockBackendRepo) UpdateContainer(containerInfo domain.ContainerInfo) error {
	return m.UpdateContainerFunc(containerInfo)
}
func (m *MockBackendRepo) PostPing(pingRes domain.PingResult) error {
	return m.PostPingFunc(pingRes)
}

type MockPingService struct {
	PingFunc func(ip string) (status string, latency float32, err error)
}
func (m *MockPingService) Ping(ip string) (string, float32, error) {
	return m.PingFunc(ip)
}
