package eventsocket

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func DialEventServer(u *url.URL, handleFunc func(io.Reader)) (io.Closer, error) {
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
	handlerFunc   func(io.Reader)
}

func newEventServerConn(c *websocket.Conn, hFunc func(io.Reader)) *eventServerConn {
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
	mType, data, err := this.websocketConn.ReadMessage()

	if err != nil {
		log.Println(err)
	} else {
		if mType == websocket.BinaryMessage {
			buffer := bytes.NewBuffer(data)
			this.handlerFunc(buffer)
		}
	}
}

func (this *eventServerConn) Close() error {
	this.doneChan <- true
	return nil
}
