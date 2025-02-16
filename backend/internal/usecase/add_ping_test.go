package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"backend/internal/domain"
	i "backend/internal/infrastructure"
	"backend/internal/usecase"
)

func TestCreatePing(t *testing.T) {
	now := time.Now()

	tests := map[string]struct {
		createFunc  func(ctx context.Context, p *domain.Ping) error
		inputPing   domain.Ping
		expectedErr error
	}{
		"successful ping creation": {
			createFunc: func(ctx context.Context, p *domain.Ping) error {
				return nil
			},
			inputPing: domain.Ping{
				IP:       "192.168.1.1",
				Latency:  10,
				PingTime: now,
			},
			expectedErr: nil,
		},
		"error in PingRepo.Create": {
			createFunc: func(ctx context.Context, p *domain.Ping) error {
				return errors.New("database error")
			},
			inputPing: domain.Ping{
				IP:       "10.0.0.1",
				Latency:  20,
				PingTime: now,
			},
			expectedErr: errors.New("database error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pingRepo := &i.MockPingRepo{
				CreateFunc: tc.createFunc,
			}

			uc := usecase.NewPingUseCase(pingRepo)
			err := uc.CreatePing(context.Background(), &tc.inputPing)

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
