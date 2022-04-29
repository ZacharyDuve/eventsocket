package eventsocket

// import (
// 	"net/http"
// 	"net/url"
// 	"testing"
// 	"time"

// 	"github.com/google/uuid"
// )

// func TestThatSimpleClientAndServerWorks(t *testing.T) {
// 	eventServer := NewEventServer()
// 	http.Handle("/", eventServer)
// 	go http.ListenAndServe(":8080", nil)
// 	var rEvent Event
// 	clientHandlerFunc := func(e Event) {
// 		t.Log("Got an event", e)
// 		rEvent = e
// 	}
// 	time.Sleep(time.Second * 1)
// 	_, err := DialEventServer(&url.URL{Host: "localhost:8080", Path: "/"}, clientHandlerFunc)
// 	if err != nil {
// 		t.Fail()
// 	}
// 	sEvent, _ := NewEvent("SomeEvent", time.Now(), uuid.New(), "I am awesome")
// 	eventServer.Send(sEvent)

// 	time.Sleep(time.Second * 1)
// 	t.Log(rEvent)

// 	if rEvent.Name() != sEvent.Name() ||
// 		rEvent.Data() != sEvent.Data() ||
// 		rEvent.OriginID().String() != sEvent.OriginID().String() ||
// 		!rEvent.OriginTime().Equal(sEvent.OriginTime()) {
// 		t.Fail()
// 	}
// }
