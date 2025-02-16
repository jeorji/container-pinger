package usecase_test

import (
	"bytes"
	"context"
	"errors"
	"log"
	"testing"

	"pinger/internal/domain"
	i "pinger/internal/infrastructure"
	"pinger/internal/usecase"
)

func TestPingContainers(t *testing.T) {
	tests := map[string]struct {
		listContainersFunc   func(all bool) ([]domain.ContainerInfo, error)
		inspectContainerFunc func(containerID string) (domain.ContainerInfo, []string, error)
		pingFunc             func(ip string) (string, float32, error)
		postPingFunc         func(pingRes domain.PingResult) error
		expectedErr          error
	}{
		"successful pings": {
			listContainersFunc: func(all bool) ([]domain.ContainerInfo, error) {
				return []domain.ContainerInfo{
					{ID: "1", Name: "container1"},
				}, nil
			},
			inspectContainerFunc: func(containerID string) (domain.ContainerInfo, []string, error) {
				return domain.ContainerInfo{ID: containerID, Name: "container" + containerID}, []string{"192.168.1.1"}, nil
			},
			pingFunc: func(ip string) (string, float32, error) {
				return "succes", 10, nil
			},
			postPingFunc: func(pingRes domain.PingResult) error {
				return nil
			},
			expectedErr: nil,
		},
		"error in ListContainers": {
			listContainersFunc: func(all bool) ([]domain.ContainerInfo, error) {
				return nil, errors.New("docker error")
			},
			inspectContainerFunc: nil,
			pingFunc:             nil,
			postPingFunc:         nil,
			expectedErr:          errors.New("docker error"),
		},
		"error in InspectContainer": {
			listContainersFunc: func(all bool) ([]domain.ContainerInfo, error) {
				return []domain.ContainerInfo{
					{ID: "1", Name: "container1"},
				}, nil
			},
			inspectContainerFunc: func(containerID string) (domain.ContainerInfo, []string, error) {
				return domain.ContainerInfo{}, nil, errors.New("inspect error")
			},
			pingFunc:     nil,
			postPingFunc: nil,
			expectedErr:  nil,
		},
		"error in PingService.Ping": {
			listContainersFunc: func(all bool) ([]domain.ContainerInfo, error) {
				return []domain.ContainerInfo{
					{ID: "1", Name: "container1"},
				}, nil
			},
			inspectContainerFunc: func(containerID string) (domain.ContainerInfo, []string, error) {
				return domain.ContainerInfo{ID: containerID, Name: "container" + containerID}, []string{"192.168.1.1"}, nil
			},
			pingFunc: func(ip string) (string, float32, error) {
				return "", 0, errors.New("ping error")
			},
			postPingFunc: func(pingRes domain.PingResult) error {
				return nil
			},
			expectedErr:  nil,
		},
		"error in BackendRepo.PostPing": {
			listContainersFunc: func(all bool) ([]domain.ContainerInfo, error) {
				return []domain.ContainerInfo{
					{ID: "1", Name: "container1"},
				}, nil
			},
			inspectContainerFunc: func(containerID string) (domain.ContainerInfo, []string, error) {
				return domain.ContainerInfo{ID: containerID, Name: "container" + containerID}, []string{"192.168.1.1"}, nil
			},
			pingFunc: func(ip string) (string, float32, error) {
				return "succes", 10, nil
			},
			postPingFunc: func(pingRes domain.PingResult) error {
				return errors.New("backend error")
			},
			expectedErr: nil,
		},
		"empty container list": {
			listContainersFunc: func(all bool) ([]domain.ContainerInfo, error) {
				return []domain.ContainerInfo{}, nil
			},
			inspectContainerFunc: nil,
			pingFunc:             nil,
			postPingFunc:         nil,
			expectedErr:          nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var logOutput bytes.Buffer
			log.SetOutput(&logOutput)
			defer log.SetOutput(nil)

			dockerRepo := &i.MockContainerRepo{
				ListContainersFunc:   tc.listContainersFunc,
				InspectContainerFunc: tc.inspectContainerFunc,
			}

			backendRepo := &i.MockBackendRepo{
				PostPingFunc: tc.postPingFunc,
			}

			pingService := &i.MockPingService{
				PingFunc: tc.pingFunc,
			}

			uc := usecase.NewPingContainersUseCase(dockerRepo, backendRepo, pingService)
			err := uc.PingContainers(context.Background(), true)

			if tc.expectedErr != nil {
				if err == nil || err.Error() != tc.expectedErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tc.expectedErr, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
