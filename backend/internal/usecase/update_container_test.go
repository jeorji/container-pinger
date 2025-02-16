package usecase_test

import (
	"context"
	"errors"
	"testing"

	"backend/internal/domain"
	i "backend/internal/infrastructure"
	"backend/internal/usecase"
)

func TestCreateOrUpdateContainerByID(t *testing.T) {
	tests := map[string]struct {
		upsertByIDFunc func(ctx context.Context, c *domain.Container) error
		inputContainer domain.Container
		expectedErr    error
	}{
		"successful update": {
			upsertByIDFunc: func(ctx context.Context, c *domain.Container) error {
				return nil
			},
			inputContainer: domain.Container{
				ID:    "1",
				Name:  "container1",
				Image: "image1",
			},
			expectedErr: nil,
		},
		"error in UpsertByID": {
			upsertByIDFunc: func(ctx context.Context, c *domain.Container) error {
				return errors.New("database error")
			},
			inputContainer: domain.Container{
				ID:    "2",
				Name:  "container2",
				Image: "image2",
			},
			expectedErr: errors.New("database error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			containerRepo := &i.MockContainerRepo{
				UpsertByIDFunc: tc.upsertByIDFunc,
			}

			uc := usecase.NewContainerUpdateUseCase(containerRepo)
			err := uc.CreateOrUpdateContainerByID(context.Background(), &tc.inputContainer)

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
