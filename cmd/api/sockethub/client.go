package sockethub

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/coder/websocket"
	"github.com/krutip7/chat-app-server/cmd/api/auth"
)

type Client struct {
	addr    string
	conn    *websocket.Conn
	msgChan chan Message
}

func (client *Client) authenticate(auth *auth.Auth) (string, bool) {
	msg, err := client.readMessage()
	if err != nil || msg.Type != "AUTH" {
		client.conn.Close(websocket.StatusGoingAway, "Unauthorized access")
		return "", false
	}

	claims, err := auth.VerifyJWT(msg.Content)
	if err != nil {
		log.Println(client.addr, "Closing connection as authentication failed with err: %v\n", err)
		client.conn.Close(websocket.StatusAbnormalClosure, "Unable to read data")
		return "", false
	}
	log.Println(client.addr, "Socket Client Authentication successful")

	userId := claims.Subject
	return userId, true
}

func (client *Client) listenAndWrite() {

	for {
		msg := <-client.msgChan

		byteStream, err := json.Marshal(msg)
		if err != nil {
			log.Printf(client.addr, "Closing connection due to error while marsahlling json: %v\n", err)
			client.conn.Close(websocket.StatusInternalError, "Unable to send data")
			break
		}

		err = client.conn.Write(context.Background(), websocket.MessageText, byteStream)
		if err != nil {
			log.Printf(client.addr, "Closing connection due to error while sending message: %v\n", err)
			client.conn.Close(websocket.StatusAbnormalClosure, "Unable to send data")
			break
		}

		log.Println(client.addr, "Sent back Message ", msg)
	}
}

func (client *Client) readMessage() (Message, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	var msg Message

	_, byteStream, err := client.conn.Read(ctx)
	if err != nil {
		log.Println(client.addr, "Closing connection due to error while receiving message: ", err)
		client.conn.Close(websocket.StatusAbnormalClosure, "Unable to read data")
		return msg, err
	}

	err = json.Unmarshal(byteStream, &msg)
	if err != nil {
		log.Println(client.addr, "Closing connection due to error while unmarshalling json: ", err)
		client.conn.Close(websocket.StatusInvalidFramePayloadData, "Unable to read data")
		return msg, err
	}

	log.Println(client.addr, "Received Message ", msg)
	return msg, nil
}