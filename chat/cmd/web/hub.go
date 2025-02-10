package main

import (
	"log"

	"github.com/jorbort/42-matcha/chat/internals/models"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	app        *aplication
}

func newHub(app *aplication) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		app:        app,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			conversationID, err := h.app.models.GetOrCreateConversation(message.From, message.To)
			if err != nil {
				log.Println(err.Error())
				return
			}
			dbMessage := &models.Message{
				ConversationID: conversationID,
				SenderID:       message.From,
				ReceiverID:     message.To,
				Message:        message.Content,
			}
			err = h.app.models.SaveMessage(dbMessage)
			if err != nil {
				log.Println(err.Error())
				return
			}
			for client := range h.clients {
				if message.To == client.userID {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
