package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Ragnar-BY/event-collector/internal/domain"
	mock_service "github.com/Ragnar-BY/event-collector/internal/mocks/repository"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSaveEvents(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	ctx := context.Background()

	testErr := errors.New("some error")
	tests := []struct {
		name        string
		events      []domain.Event
		err         error
		expectedErr error
	}{
		{
			name: "success case",
			events: []domain.Event{
				{
					ClientTime: time.Now(),
					Sequence:   uint64(1),
					DeviceID:   uuid.New(),
					Session:    "session",
				}},
			err:         nil,
			expectedErr: nil,
		},
		{
			name:        "error case",
			events:      nil,
			err:         testErr,
			expectedErr: testErr,
		},
	}

	MockRepository := mock_service.NewMockRepository(ctl)
	repo := NewRepositoryService(MockRepository)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gomock.InOrder(
				MockRepository.EXPECT().SaveEvents(ctx, test.events).Return(test.err),
			)
			err := repo.SaveEvents(ctx, test.events)
			if !assert.Equal(t, err, test.expectedErr) {
				t.Errorf("repo.SaveEvents: expected %v, got %v", test.expectedErr, err)
			}
		})
	}
}
