package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ClientTime time.Time
	DeviceID   uuid.UUID
	DeviceOS   string
	Session    string
	Sequence   uint64
	Event      string
	ParamInt   uint64
	ParamStr   string
	ClientIP string
	ServerTime time.Time
}
