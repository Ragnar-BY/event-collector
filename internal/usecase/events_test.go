package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Ragnar-BY/event-collector/internal/domain"
	mock_usecase "github.com/Ragnar-BY/event-collector/internal/mocks/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSaveEvents(t *testing.T) {

	clientIP := "11.22.33.44"
	serverTime := time.Now()
	logger, _ := zap.NewDevelopment()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	eventService := mock_usecase.NewMockEventService(ctl)

	ctx := context.Background()

	events := []domain.Event{}
	eventsJson := []byte(`[
	{"ClientTime":"2020-12-01T23:59:00Z","DeviceID":"0287d9aa-4adf-4b37-a60f-3e9e645c821e","DeviceOS":"iOS 13.5.1","Session":"ybuRi8mAUypxjbxQ","Sequence":1,"Event":"app_start","ParamInt":0,"ParamStr":"some text","ClientIP":"","ServerTime":"0001-01-01T00:00:00Z"},
	{"ClientTime":"2020-12-01T23:59:00Z","DeviceID":"0287d9aa-4adf-4b37-a60f-3e9e645c821e","DeviceOS":"iOS 13.5.1","Session":"ybuRi8mAUypxjbxQ","Sequence":2,"Event":"app_start","ParamInt":0,"ParamStr":"some text","ClientIP":"","ServerTime":"0001-01-01T00:00:00Z"}
	]`)
	err := json.Unmarshal(eventsJson, &events)
	if err != nil {
		t.Error(err)

		return
	}
	for i := range events {
		events[i].ClientIP = clientIP
		events[i].ServerTime = serverTime
	}

	gomock.InOrder(
		eventService.EXPECT().SaveEvents(ctx, events).Return(nil),
	)

	use := NewEventUsecase(eventService, logger, Config{NumberOfThreads: 1, ChannelCapacity: 1})

	err = use.SaveEvents(ctx, events, clientIP, serverTime)
	time.Sleep(1 * time.Millisecond) // wait to send batch to channel

	assert.Equal(t, nil, err)
}
