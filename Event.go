package eventsocket

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Event interface {
	Name() string
	OriginTime() time.Time
	OriginID() uuid.UUID
	Data() string
}

type event struct {
	name       string
	data       string
	originTime time.Time
	originID   uuid.UUID
}

func NewEvent(eName string, originID uuid.UUID, data string) (Event, error) {
	if eName == "" {
		return nil, errors.New("eName is required for event")
	}

	return &event{name: eName, data: data, originTime: time.Now(), originID: originID}, nil
}

func (this *event) Name() string {
	return this.name
}

func (this *event) OriginTime() time.Time {
	return this.originTime
}
func (this *event) OriginID() uuid.UUID {
	return this.originID
}
func (this *event) Data() string {
	return this.data
}
