package usecase_test

import (
	"context"
	"errors"
	"reflect"

	"testing"
	"time"

	"backend/internal/domain"
	i "backend/internal/infrastructure"
	"backend/internal/usecase"
)

func TestGetAllStats(t *testing.T) {
	now := time.Now()

	tests := map[string]struct {
		getAllFunc              func(ctx context.Context) ([]domain.Container, error)
		getLastPingForContainer func(ctx context.Context, container *domain.Container) ([]domain.Ping, error)
		expectedResult          []domain.ContainerStatsDTO
		expectedErr             error
	}{
		"error in ContainerRepo.GetAll": {
			getAllFunc: func(ctx context.Context) ([]domain.Container, error) {
				return nil, errors.New("container error")
			},
			getLastPingForContainer: func(ctx context.Context, container *domain.Container) ([]domain.Ping, error) {
				return []domain.Ping{}, nil
			},
			expectedResult: nil,
			expectedErr:    errors.New("container error"),
		},
		"error in PingRepo.GetLastPingForContainer": {
			getAllFunc: func(ctx context.Context) ([]domain.Container, error) {
				return []domain.Container{
					{
						ID:        "1",
						Name:      "container1",
						Image:     "image1",
						State:     "running",
						Status:    "ok",
						CreatedAt: now,
					},
				}, nil
			},
			getLastPingForContainer: func(ctx context.Context, container *domain.Container) ([]domain.Ping, error) {
				return nil, errors.New("ping error")
			},
			expectedResult: nil,
			expectedErr:    errors.New("ping error"),
		},
		"successful execution": {
			getAllFunc: func(ctx context.Context) ([]domain.Container, error) {
				return []domain.Container{
					{
						ID:        "1",
						Name:      "container1",
						Image:     "image1",
						State:     "running",
						Status:    "ok",
						CreatedAt: now,
					},
					{
						ID:        "2",
						Name:      "container2",
						Image:     "image2",
						State:     "stopped",
						Status:    "not ok",
						CreatedAt: now,
					},
				}, nil
			},
			getLastPingForContainer: func(ctx context.Context, container *domain.Container) ([]domain.Ping, error) {
				switch container.ID {
				case "1":
					return []domain.Ping{
						{IP: "172.17.0.2", Latency: 10, PingTime: now},
					}, nil
				case "2":
					return []domain.Ping{
						{IP: "172.17.0.3", Latency: 20, PingTime: now},
						{IP: "172.19.0.2", Latency: 30, PingTime: now},
					}, nil
				default:
					return nil, nil
				}
			},
			expectedResult: []domain.ContainerStatsDTO{
				{
					ID:        "1",
					Name:      "container1",
					Image:     "image1",
					State:     "running",
					Status:    "ok",
					CreatedAt: now,
					PingStats: []domain.PingStatsDTO{
						{IP: "172.17.0.2", LastLatency: 10, LastSuccessful: now},
					},
				},
				{
					ID:        "2",
					Name:      "container2",
					Image:     "image2",
					State:     "stopped",
					Status:    "not ok",
					CreatedAt: now,
					PingStats: []domain.PingStatsDTO{
						{IP: "172.17.0.3", LastLatency: 20, LastSuccessful: now},
						{IP: "172.19.0.2", LastLatency: 30, LastSuccessful: now},
					},
				},
			},
			expectedErr: nil,
		},
		"empty container list": {
			getAllFunc: func(ctx context.Context) ([]domain.Container, error) {
				return []domain.Container{}, nil
			},
			getLastPingForContainer: func(ctx context.Context, container *domain.Container) ([]domain.Ping, error) {
				return []domain.Ping{}, nil
			},
			expectedResult: []domain.ContainerStatsDTO{},
			expectedErr:    nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			containerRepo := &i.MockContainerRepo{
				GetAllFunc: tc.getAllFunc,
			}
			pingRepo := &i.MockPingRepo{
				GetLastPingForContainerFunc: tc.getLastPingForContainer,
			}

			uc := usecase.NewContainerStatsUseCase(containerRepo, pingRepo)
			stats, err := uc.GetAllStats(context.Background())

			if tc.expectedErr != nil {
				if err == nil || err.Error() != tc.expectedErr.Error() {
					t.Fatalf("expected error: %v, got: %v", tc.expectedErr, err)
				}
			}

			if !reflect.DeepEqual(stats, tc.expectedResult) {
				t.Errorf("expected: %+v, got: %+v", tc.expectedResult, stats)
			}
		})
	}
}
