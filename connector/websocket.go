package connector

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/cihub/seelog"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsAccessToken = "";

func SetAccessToken(accessToken string) {
	seelog.Infof("access token : %s", accessToken)
	wsAccessToken = accessToken
}

func ServeWs(hub *WsGroup, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		seelog.Error(err)
		return
	}

	r.ParseForm()
	accessToken := r.Header.Get("access_token")
	if accessToken == "" {
		accessToken = r.Form.Get("access_token")
	}

	if accessToken != wsAccessToken {
		conn.Close()
		seelog.Errorf("WS access_token error: %s", accessToken)
		return
	}

	client := &WsClient{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	hub.register <- client
	go client.push()
	client.listen()
}
