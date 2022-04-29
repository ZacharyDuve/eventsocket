package eventsocket

import (
	"errors"
	"io"
	"log"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func DialEventServer(u *url.URL, handleFunc func(Event)) (io.Closer, error) {
	if u == nil {
		return nil, errors.New("u of URL is required to dial an event server")
	}
	if handleFunc == nil {
		return nil, errors.New("func(Event) is required to dial an event server")
	}
	if u.Scheme != "ws" {
		u.Scheme = "ws"
	}
	var closer io.Closer
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err == nil {
		closer = newEventServerConn(c, handleFunc)
	}

	return closer, err
}

type eventServerConn struct {
	websocketConn *websocket.Conn
	doneChan      chan bool
	handlerFunc   func(Event)
}

func newEventServerConn(c *websocket.Conn, hFunc func(Event)) *eventServerConn {
	eSConn := &eventServerConn{websocketConn: c}
	eSConn.doneChan = make(chan bool)
	eSConn.handlerFunc = hFunc
	go eSConn.processEventsLoop()
	return eSConn
}

func (this *eventServerConn) processEventsLoop() {
	for {
		select {
		case <-this.doneChan:
			this.websocketConn.Close()
			break
		default:
			this.processEvent()
		}
	}
}

func (this *eventServerConn) processEvent() {
	jEvent := &jsonEvent{}
	err := this.websocketConn.ReadJSON(jEvent)
	if err == nil {
		var id uuid.UUID
		id, err = uuid.Parse(jEvent.OriginUUID)
		if err == nil {
			var e Event
			e, err = NewEvent(jEvent.EventName, jEvent.OriginTime, id, jEvent.Data)
			if err == nil {
				this.handlerFunc(e)
			}
		}
	}

	if err != nil {
		log.Println(err)
	}
}

func (this *eventServerConn) Close() error {
	this.doneChan <- true
	return nil
}
