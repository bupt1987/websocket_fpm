package connector

import (
	"runtime"
	"github.com/cihub/seelog"
	"github.com/bupt1987/log-websocket/util"
)

type WsGroup struct {
	num        int64
	clients    map[*WsClient]bool
	Broadcast  chan []byte
	register   chan *WsClient
	unregister chan *WsClient
}

func NewWsGroup() *WsGroup {
	return &WsGroup{
		num:        0,
		Broadcast:  make(chan []byte, runtime.NumCPU()),
		register:   make(chan *WsClient),
		unregister: make(chan *WsClient),
		clients:    make(map[*WsClient]bool),
	}
}

func (h *WsGroup) push(client *WsClient, msg []byte) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("WsGroup.push error, ", err);
		}
	}()
	client.send <- msg
}

func (h *WsGroup) Run() {
	go func() {
		defer util.PanicExit()
		for {
			select {
			case client := <-h.register:
				h.clients[client] = true
				h.num ++;
				seelog.Debugf("%s connected, total connected: %v", client.conn.RemoteAddr(), h.num)
			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
					h.num --;
					seelog.Debugf("%s close, total connected: %v", client.conn.RemoteAddr(), h.num)
				}
			case msg := <-h.Broadcast:
				go func() {
					for client := range h.clients {
						h.push(client, msg)
					}
				}()
			}
		}
	}()
}
