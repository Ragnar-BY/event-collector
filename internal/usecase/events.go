package usecase

import (
	"context"
	"time"

	"github.com/Ragnar-BY/event-collector/internal/domain"

	"go.uber.org/zap"
)

type Config struct {
	NumberOfThreads int
	ChannelCapacity int
}

// EventService describes events service
type EventService interface {
	SaveEvents(ctx context.Context, events []domain.Event) error

	PingDB(ctx context.Context) error
	CloseDB() error
}

// EventUsecase is usecase for events
type EventUsecase struct {
	service   EventService
	log       *zap.Logger
	eventChan chan []domain.Event
}

// NewEventUsecase creates new event usecase
func NewEventUsecase(events EventService, log *zap.Logger, cfg Config) *EventUsecase {
	u := &EventUsecase{
		service:   events,
		log:       log,
		eventChan: make(chan []domain.Event, cfg.ChannelCapacity),
	}

	for i := 1; i <= cfg.NumberOfThreads; i++ {
		go u.saveEventsDaemon()
	}

	return u
}

// SaveEvents add info about clientIP and server time and send events to channel for saving
func (u *EventUsecase) SaveEvents(ctx context.Context, events []domain.Event, clientIP string, serverTime time.Time) error {
	for i := range events {
		events[i].ClientIP = clientIP
		events[i].ServerTime = serverTime
	}
	select {
	case u.eventChan <- events:
		return nil
	default:
		return domain.ErrToManyRequests
	}
}

func (u *EventUsecase) saveEventsDaemon() {
	for batch := range u.eventChan {
		err := u.service.SaveEvents(context.Background(), batch)
		if err != nil {
			u.log.Error("can not save events", zap.Error(err))
		}
	}
}

// Close closes repository
func (u *EventUsecase) Close() error {
	return u.service.CloseDB()
}

// Ping pings repository
func (u *EventUsecase) Ping(ctx context.Context) error {
	return u.service.PingDB(ctx)
}
