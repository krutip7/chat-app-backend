package sockethub

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/krutip7/chat-app-server/cmd/api/auth"
)

type WebSocketHandler struct {
	Auth    *auth.Auth
	ChatHub *Hub
}

func NewWebSocketHandler(auth *auth.Auth) *WebSocketHandler {
	return &WebSocketHandler{
		Auth:    auth,
		ChatHub: NewHub(),
	}
}

func (wsh *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {

	opts := &websocket.AcceptOptions{OriginPatterns: []string{"*"}}
	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		log.Println(r.RemoteAddr, "Could not connect due to error: ", err)
		return
	}
	log.Println(r.RemoteAddr, "Accepted connection from ", r.RemoteAddr)
	defer conn.CloseNow()

	client := Client{
		addr:    r.RemoteAddr,
		conn:    conn,
		msgChan: make(chan Message),
	}

	userId, ok := client.authenticate(wsh.Auth)
	if !ok {
		return
	}

	wsh.ChatHub.addClient(userId, client)

	go client.listenAndWrite()

	for {
		log.Println(r.RemoteAddr, "Awaiting Message")
		msg, err := client.readMessage()
		if err != nil {
			break
		}
		wsh.ChatHub.processMessage(msg)
	}

	wsh.ChatHub.removeClient(userId, client)

	log.Println(r.RemoteAddr, "Closing Connection")
	conn.Close(websocket.StatusNormalClosure, "")
}
