package eventsocket

import (
	"time"
)

type jsonEvent struct {
	EventName  string    `json:"event-name"`
	OriginTime time.Time `json:"origin-time"`
	OriginUUID string    `json:"origin-uuid"`
	Data       string    `json:"data"`
}
