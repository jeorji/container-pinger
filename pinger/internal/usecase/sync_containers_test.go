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

func TestSyncContainers(t *testing.T) {
	tests := map[string]struct {
		listContainersFunc   func(all bool) ([]domain.ContainerInfo, error)
		inspectContainerFunc func(containerID string) (domain.ContainerInfo, []string, error)
		updateContainerFunc  func(containerInfo domain.ContainerInfo) error
		expectedErr          error
	}{
		"successful sync": {
			listContainersFunc: func(all bool) ([]domain.ContainerInfo, error) {
				return []domain.ContainerInfo{
					{ID: "1", Name: "container1"},
					{ID: "2", Name: "container2"},
				}, nil
			},
			inspectContainerFunc: func(containerID string) (domain.ContainerInfo, []string, error) {
				return domain.ContainerInfo{ID: containerID, Name: "container" + containerID}, nil, nil
			},
			updateContainerFunc: func(containerInfo domain.ContainerInfo) error {
				return nil
			},
			expectedErr: nil,
		},
		"error in ListContainers": {
			listContainersFunc: func(all bool) ([]domain.ContainerInfo, error) {
				return nil, errors.New("docker error")
			},
			inspectContainerFunc: nil,
			updateContainerFunc:  nil,
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
			updateContainerFunc: nil,
			expectedErr:         nil,
		},
		"error in UpdateContainer": {
			listContainersFunc: func(all bool) ([]domain.ContainerInfo, error) {
				return []domain.ContainerInfo{
					{ID: "1", Name: "container1"},
				}, nil
			},
			inspectContainerFunc: func(containerID string) (domain.ContainerInfo, []string, error) {
				return domain.ContainerInfo{ID: containerID, Name: "container" + containerID}, nil, nil
			},
			updateContainerFunc: func(containerInfo domain.ContainerInfo) error {
				return errors.New("update error")
			},
			expectedErr: nil,
		},
		"empty container list": {
			listContainersFunc: func(all bool) ([]domain.ContainerInfo, error) {
				return []domain.ContainerInfo{}, nil
			},
			inspectContainerFunc: nil,
			updateContainerFunc:  nil,
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
				UpdateContainerFunc: tc.updateContainerFunc,
			}

			uc := usecase.NewContainerSyncUseCase(dockerRepo, backendRepo)
			err := uc.SyncContainers(context.Background(), true)

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
