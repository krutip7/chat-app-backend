package sockethub

import (
	"log"
	"slices"
)

type Hub struct {
	activeUsers map[string][]Client // userId -> List of active Client Instances
	rooms       map[string]string   // roomId -> userId
}

func NewHub() *Hub {
	return &Hub{
		activeUsers: make(map[string][]Client),
		rooms:       make(map[string]string),
	}
}

func (hub *Hub) addClient(userId string, client Client) {

	activeClients, ok := hub.activeUsers[userId]
	if !ok {
		activeClients = make([]Client, 0)
	}

	hub.activeUsers[userId] = append(activeClients, client)
}

func (hub *Hub) removeClient(userId string, client Client) {

	activeClients, ok := hub.activeUsers[userId]
	if !ok || activeClients == nil {
		return
	}

	activeClients = slices.DeleteFunc(activeClients, func(c Client) bool { return c.conn == client.conn })
	if len(activeClients) == 0 {
		log.Println("No more active connections for user ", userId)
		delete(hub.activeUsers, userId)
	}

}


func (hub *Hub) processMessage(msg Message) {

	if userClients, ok := hub.activeUsers[msg.From]; ok {
		for _, recipient := range userClients {
			recipient.msgChan <- msg
		}
	}

	if msg.To == msg.From {
		return
	}

	if userClients, ok := hub.activeUsers[msg.To]; ok {
		for _, recipient := range userClients {
			recipient.msgChan <- msg
		}
	}
}
