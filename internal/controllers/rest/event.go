package rest

import (
	"strings"
	"time"

	"github.com/Ragnar-BY/event-collector/internal/domain"

	"github.com/google/uuid"
)

// Event is struct for event from rest endpoint
type Event struct {
	ClientTime CustomTime `json:"client_time"`
	DeviceID   uuid.UUID  `json:"device_id"`
	DeviceOS   string     `json:"device_os"`
	Session    string     `json:"session"`
	Sequence   uint64     `json:"sequence"`
	Event      string     `json:"event"`
	ParamInt   uint64     `json:"param_int"`
	ParamStr   string     `json:"param_str"`
}

func (e Event) ToDomainEvent() domain.Event {
	return domain.Event{
		ClientTime: time.Time(e.ClientTime),
		DeviceID:   e.DeviceID,
		DeviceOS:   e.DeviceOS,
		Session:    e.Session,
		Sequence:   e.Sequence,
		Event:      e.Event,
		ParamInt:   e.ParamInt,
		ParamStr:   e.ParamStr,
	}
}

// CustomeTime is type for working with time format from client
type CustomTime time.Time

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	layout := "2006-01-02 15:04:05"

	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(layout, s)
	*ct = CustomTime(nt)
	return
}
