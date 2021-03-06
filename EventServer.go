package eventsocket

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type EventServer interface {
	Send(io.Reader)
	http.Handler
	io.Closer
}

type eventServer struct {
	upgrader     websocket.Upgrader
	clients      []*websocket.Conn
	clientsMutex *sync.Mutex
}

func NewEventServer() EventServer {
	eS := &eventServer{}
	eS.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	eS.clients = make([]*websocket.Conn, 0)
	eS.clientsMutex = &sync.Mutex{}
	return eS
}

func (this *eventServer) Close() error {
	this.clientsMutex.Lock()
	for _, curClient := range this.clients {
		curClient.Close()
	}
	this.clientsMutex.Unlock()
	return nil
}
func (this *eventServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := this.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	this.clientsMutex.Lock()
	this.clients = append(this.clients, c)
	this.clientsMutex.Unlock()
}

func (this *eventServer) Send(r io.Reader) {
	if r != nil {
		var deadClients []*websocket.Conn
		this.clientsMutex.Lock()
		for _, curClient := range this.clients {
			data, err := ioutil.ReadAll(r)
			if err == nil {
				curClient.WriteMessage(websocket.BinaryMessage, data)
			} else {
				log.Println(err)
				if websocket.IsUnexpectedCloseError(err) {
					if deadClients == nil {
						deadClients = make([]*websocket.Conn, 0)
					}
					deadClients = append(deadClients, curClient)
				}
			}
		}
		if len(deadClients) > 0 {
			this.removeDeadClients(deadClients)
		}
		this.clientsMutex.Unlock()
	}
}

func (this *eventServer) removeDeadClients(dClients []*websocket.Conn) {
	this.clientsMutex.Lock()
	newClientsSlice := make([]*websocket.Conn, 0, len(this.clients)-len(dClients))
	for _, curClient := range this.clients {
		isDead := false
		for _, curDeadClient := range dClients {
			if curDeadClient == curClient {
				isDead = true
				break
			}
		}
		if !isDead {
			newClientsSlice = append(newClientsSlice, curClient)
		}
	}

	this.clientsMutex.Unlock()
}
